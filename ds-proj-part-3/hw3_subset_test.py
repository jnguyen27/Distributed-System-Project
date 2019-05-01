import os
import sys
import requests
import time
import unittest
import json

import docker_control

import io

dockerBuildTag = "grslin/gjjr_kvs" #put the tag for your docker build here

hostIp = "localhost" # this can be localhost again

needSudo = False # obviously if you need sudo, set this to True 
#contact me imediately if setting this to True breaks things 
#(I don't have a machine which needs sudo, so it has not been tested, although in theory it should be fine)

port_prefix = "808"

networkName = "mynet" # the name of the network you created

networkIpPrefix = "192.168.1." # should be everything up to the last period of the subnet you specified when you 
# created your network

propogationTime = 3 #sets number of seconds we sleep after certain actions to let data propagate through your system
# you may lower this to speed up your testing if you know that your system is fast enough to propigate information faster than this
# I do not recomend increasing this

shouldJsonify = True

dc = docker_control.docker_controller(networkName, needSudo)

def jsonify(payload):
    if shouldJsonify:
        return json.dumps(payload)
    else:
        return payload

######################################################################
## HTTP Requests:
##############################
# Basic Functionality
# These are the endpoints we should be able to hit
    #KVS Functions
def storeKeyValue(ipPort, key, value, payload):
    print(' PUT: http://%s/keyValue-store/%s'%(str(ipPort), key))
    return requests.put( 'http://%s/keyValue-store/%s'%(str(ipPort), key), data={'val':value, 'payload': jsonify(payload)})

def checkKey(ipPort, key, payload):
    print(' GET: http://%s/keyValue-store/search/%s'%(str(ipPort), key))
    return requests.get( 'http://%s/keyValue-store/search/%s'%(str(ipPort), key), data={'payload': jsonify(payload)} )

def getKeyValue(ipPort, key, payload):
    print(' GET: http://%s/keyValue-store/%s'%(str(ipPort), key))
    return requests.get( 'http://%s/keyValue-store/%s'%(str(ipPort), key), data={'payload': jsonify(payload)} )

def deleteKey(ipPort, key, payload):
    print(' DELETE: http://%s/keyValue-store/%s'%(str(ipPort), key))
    return requests.delete( 'http://%s/keyValue-store/%s'%(str(ipPort), key), data={'payload': jsonify(payload)} )

    #Replication Functions
def addNode(ipPort, newAddress):
    print(' PUT: http://%s/view'%str(ipPort))
    return requests.put( 'http://%s/view'%str(ipPort), data={'ip_port':newAddress} )

def removeNode(ipPort, oldAddress):
    print(' DELETE: http://%s/view'%str(ipPort))
    return requests.delete( 'http://%s/view'%str(ipPort), data={'ip_port':oldAddress} )

def viewNetwork(ipPort):
    print(' GET: http://%s/view'%str(ipPort))
    return requests.get( 'http://%s/view'%str(ipPort) )

###########################################################################################

class TestHW3Subset(unittest.TestCase):
    view = {}

    def setUp(self):
        dc.cleanUpDockerContainer()
        self.view = dc.spinUpManyContainers(dockerBuildTag, hostIp, networkIpPrefix, port_prefix, 2)

        for container in self.view:
            if " " in container["containerID"]:
                self.assertTrue(False, "There is likely a problem in the settings of your ip addresses or network.")

        dc.prepBlockade([instance["containerID"] for instance in self.view])

    def tearDown(self):
        dc.cleanUpDockerContainer()
        dc.tearDownBlockade()

##############################################################################################
## Internal Helpers
#############################################
    def getPayload(self, ipPort, key):
        response = checkKey(ipPort, key, {})
        #print(response)
        data = response.json()
        return data["payload"]

    def partition(self, partitionList):
        truncatedList = []
        for partition in partitionList:
            truncatedPartition = []
            for node in partition:
                truncatedPartition.append(node[:12])
            truncatedPartition = ",".join(truncatedPartition)
            truncatedList.append(truncatedPartition)
        dc.partitionContainer(truncatedList)

    def confirmAddKey(self, ipPort, key, value, expectedStatus, expectedMsg, expectedReplaced, payload={}):
        response = storeKeyValue(ipPort, key, value, payload)
        print("     Status Code: %s"%response.status_code)
        print("     Expected Status Code: %s"%expectedStatus)

        self.assertEqual(response.status_code, expectedStatus)

        data = response.json()
        print (data)
        
        print("     Response: ")
        print("         Message: %s"%data['msg'])
        print("         Replaced: %s"%data['replaced'])
        print("         payload: %s"%data['payload'])
        print("     Expected: ")
        print("         Message: %s"%expectedMsg)
        print("         Replaced: %s"%expectedReplaced)

        self.assertEqual(data['msg'], expectedMsg)
        self.assertEqual(data['replaced'], expectedReplaced)

        return data["payload"]

    def confirmCheckKey(self, ipPort, key, expectedStatus, expectedResult, expectedIsExists, payload={}):
        response = checkKey(ipPort, key, payload)
        print("     Status Code: %s"%response.status_code)
        print("     Expected Status Code: %s"%expectedStatus)
        self.assertEqual(response.status_code, expectedStatus)

        data = response.json()

        print("     Response: ")
        print("         Result: %s"%data['result'])
        print("         isExists: %s"%data['isExists'])
        print("         payload: %s"%data['payload'])
        print("     Expected: ")
        print("         Result: %s"%expectedResult)
        print("         isExists: %s"%expectedIsExists)

        self.assertEqual(data['result'], expectedResult)
        self.assertEqual(data['isExists'], expectedIsExists)

        return data["payload"]

    def confirmGetKey(self, ipPort, key, expectedStatus, expectedResult, expectedValue=None, expectedMsg=None, payload={}):
        response = getKeyValue(ipPort, key, payload)
        print("     Status Code: %s"%response.status_code)
        print("     Expected Status Code: %s"%expectedStatus)
        self.assertEqual(response.status_code, expectedStatus)

        data = response.json()

        print("     Response: ")
        if expectedMsg != None and 'msg' in data:
            print("         Message: %s"%data['msg'])
        print("         Result: %s"%data['result'])
        if expectedValue != None and 'value' in data:
            print("         Value: %s"%data['value'])
        print("         payload: %s"%data['payload'])
        print("     Expected: ")
        if expectedMsg != None and 'msg' in data:
            print("         Message: %s"%expectedMsg)
        print("         Result: %s"%expectedResult)
        if expectedValue != None and 'value' in data:
            print("         Value: %s"%expectedValue)

        self.assertEqual(data['result'], expectedResult)
        if expectedValue != None and 'value' in data:
            self.assertEqual(data['value'], expectedValue)
        if expectedMsg != None and 'msg' in data:
            self.assertEqual(data['msg'], expectedMsg)

        return data["payload"]



#########################################################################################################################
## Actual Tests:
#######################################

    def test_q_causual_consistency(self):
      print("begining test Q:")
      keyX = "TheKeyX"
      valueZero = "Zero"

      nodeA = self.view[0]["testScriptAddress"]
      nodeB = self.view[1]["testScriptAddress"]
      
      print("An initial client writes X=0 (no causal history) to A")
      self.confirmAddKey(ipPort = nodeA, 
                         key=keyX, 
                         value=valueZero, 
                         expectedStatus=200, 
                         expectedMsg="Added successfully", 
                         expectedReplaced=False)

      print("Client X reads X (=0) from A, writes X=1 to B (with appropriate causal history).")
      payloadClient1 = self.confirmGetKey(ipPort=nodeA, 
                           key=keyX, 
                           expectedStatus=200, 
                           expectedResult="Success", 
                           expectedValue=valueZero)

      valueOne = "One"

      time.sleep(propogationTime)

      payloadClient1 = self.confirmAddKey(ipPort = nodeB, 
                         key=keyX, 
                         value=valueOne, 
                         expectedStatus=201, 
                         expectedMsg="Updated successfully", 
                         expectedReplaced=True,
                         payload= payloadClient1)

      time.sleep(propogationTime)

      print("Client Y then reads X=1, writes Y=1 (both to B).")
      payloadClient2 = self.confirmGetKey(ipPort=nodeB, 
                           key=keyX, 
                           expectedStatus=200, 
                           expectedResult="Success", 
                           expectedValue=valueOne)

      keyY = "TheKeyY"

   

      payloadClient2 = self.confirmAddKey(ipPort = nodeB, 
                         key=keyY, 
                         value=valueOne, 
                         expectedStatus=200, 
                         expectedMsg="Added successfully", 
                         expectedReplaced=False,
                         payload= payloadClient1)

      time.sleep(propogationTime)

      print("Client X then reads Y=1 from B and reads X=1 from A.")
      payloadClient1 = self.confirmGetKey(ipPort=nodeB, 
                           key=keyY, 
                           expectedStatus=200, 
                           expectedResult="Success", 
                           expectedValue=valueOne,
                           payload=payloadClient1)

      payloadClient1 = self.confirmGetKey(ipPort=nodeA, 
                           key=keyX, 
                           expectedStatus=200, 
                           expectedResult="Success", 
                           expectedValue=valueOne,
                           payload=payloadClient1)

    def test_r_causual_consistency_partition(self):
        print("begining test R:")
        keyX = "TheKeyX"
        valueZero = "Zero"

        nodeA = self.view[0]["testScriptAddress"]
        nodeAContainerID = self.view[0]["containerID"]

        nodeB = self.view[1]["testScriptAddress"]
        nodeBContainerID = self.view[1]["containerID"]

        print("partition A and B")
        self.partition([[nodeAContainerID], [nodeBContainerID]])

        print("An initial client writes X=0 (no causal history) to A")
        self.confirmAddKey(ipPort = nodeA, 
                         key=keyX, 
                         value=valueZero, 
                         expectedStatus=200, 
                         expectedMsg="Added successfully", 
                         expectedReplaced=False)

        print("Client X reads X (=0) from A, writes X=1 to B (with appropriate causal history).")
        payloadClient1 = self.confirmGetKey(ipPort=nodeA, 
                           key=keyX, 
                           expectedStatus=200, 
                           expectedResult="Success", 
                           expectedValue=valueZero)

        valueOne = "One"

        payloadClient1 = self.confirmAddKey(ipPort = nodeB, 
                         key=keyX, 
                         value=valueOne, 
                         expectedStatus=200, 
                         expectedMsg="Added successfully", 
                         expectedReplaced=False,
                         payload= payloadClient1)

        print("Client Y then reads X=1, writes Y=1 (both to B).")
        payloadClient2 = self.confirmGetKey(ipPort=nodeB, 
                           key=keyX, 
                           expectedStatus=200, 
                           expectedResult="Success", 
                           expectedValue=valueOne)

        keyY = "TheKeyY"

        payloadClient2 = self.confirmAddKey(ipPort = nodeB, 
                         key=keyY, 
                         value=valueOne, 
                         expectedStatus=200, 
                         expectedMsg="Added successfully", 
                         expectedReplaced=False,
                         payload= payloadClient1)

        print("Client X then reads Y=1 from B and attempts to read X from A.")
        payloadClient1 = self.confirmGetKey(ipPort=nodeB, 
                           key=keyY, 
                           expectedStatus=200, 
                           expectedResult="Success", 
                           expectedValue=valueOne,
                           payload=payloadClient1)

        print("Client X should NOT return X=0!")

        response = getKeyValue(nodeA, keyX, payloadClient1)

        print("     Status Code: %s"%response.status_code)
        print("     Expected Status Code: 400, or less than 300")


        if response.status_code < 300:
            response = response.json()

            print("     Response: ")
            print("         Value: %s"%response['value'])
            print("         payload: %s"%response['payload'])
            print("     Expected: ")
            print("         Value: %s"%valueZero)

            self.assertNotEqual(response["value"], valueZero)
        else: 
            self.assertEqual(response.status_code, 400)



if __name__ == '__main__':

    unittest.main() 
