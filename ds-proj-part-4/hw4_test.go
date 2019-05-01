/**********************************************
/ AUTHOR: Akobir Khamidov (akhamido@ucsc.edu) *
/ AUTHOR: Vien Van (vhvan@ucsc.edu)			  *
/ COPYRIGHT 2018 © by TEAMAWESOME			  *
***********************************************/
package hw4

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestShardIdIsUniqueFromEveryNode(t *testing.T) {
	runContainers(2, 4)
	key := "GuessWhat??"
	value := "Vien is Awesome!"
	initPutKey(t, key, value)

	var ownerIDs []float64
	for port := range containersInfos {
		res := sendRequest(port, "GET", "keyValue-store/"+key, t, http.StatusOK, nil)
		owner, _ := res["owner"]
		ownerID := owner.(float64)
		ownerIDs = append(ownerIDs, ownerID)
	}
	for i := 1; i < len(ownerIDs); i++ {
		if ownerIDs[i] != ownerIDs[0] {
			t.Errorf("Shard IDs are not unique!")
		}
	}
}

func TestDeleteNodeAllKeysShouldMigrate(t *testing.T) {
	runContainers(2, 4)
	miniKVS := map[string]int{
		"A": 1, "B": 2, "C": 3,
		"D": 4, "E": 5, "F": 6,
		"G": 7, "H": 8, "I": 9,
		"J": 10, "K": 11, "L": 12,
	}
	for k, v := range miniKVS {
		initPutKey(t, k, string(v))
	}
	time.Sleep(2 * time.Second)
	// for k, v := range miniKVS {
	// 	expectedGetKey := map[string]interface{}{"result": "Success", "value": v, "statusCode": http.StatusOK}
	// 	ownerID := checkGetKey(t, k, expectedGetKey, !askEveryNode)
	// 	log.Println("OwnerID: ", ownerID)
	// }

	node := getPort()
	log.Println("Node: ", node)
	deleteNode(t, node, !askEveryNode)
	time.Sleep(2 * time.Second)
	count := getCount(t, 0, http.StatusOK, !askEveryNode)
	log.Println("Count: ", count)
	// for k, v := range miniKVS {
	// 	expectedGetKey := map[string]interface{}{"result": "Success", "value": v, "statusCode": http.StatusOK}
	// 	ownerID := checkGetKey(t, k, expectedGetKey, !askEveryNode)
	// 	// count := getCount(t, ownerID, http.StatusOK, !askEveryNode)
	// 	if ownerID != 0 {
	// 		t.Errorf("Keys did not migrate. Owner ID %v still exists\n", ownerID)
	// 	}
	// 	log.Println("OwnerID: ", ownerID)
	// }

}

func DeleteKeyShouldDecrementCount(t *testing.T) {
	runContainers(2, 4)
	key := "Hey!"
	value := "Vien is awesome!"
	initPutKey(t, key, value)
	expectedGetKey := map[string]interface{}{"result": "Success", "value": value, "statusCode": http.StatusOK}
	ownerID := checkGetKey(t, key, expectedGetKey, !askEveryNode)
	count := getCount(t, ownerID, http.StatusOK, !askEveryNode)

	expectedDeleteKey := map[string]interface{}{"result": "Success", "msg": "Key deleted", "statusCode": http.StatusOK}
	checkDeleteKey(t, key, expectedDeleteKey, !askEveryNode)
	time.Sleep(2 * time.Second)

	expectedGetKey = map[string]interface{}{"result": "Error", "msg": "Key does not exist", "statusCode": http.StatusNotFound}

	otherCount := getCount(t, ownerID, http.StatusOK, askEveryNode)

	if otherCount == count {
		t.Errorf("count before delete: %v == count after delete: %v", count, otherCount)
	}

}

func TestGetShardIDFromEveryNode(t *testing.T) {
	runContainers(2, 4)
	for port := range containersInfos {
		res := sendRequest(port, "GET", "shard/my_id", t, http.StatusOK, nil)
		if res["id"].(float64) >= float64(initNumOfShards) {
			t.Errorf("localhost:%v: shard_id must be < %v: %v", port, initNumOfShards, res["id"].(int))
		}
	}
}

func TestGetAllShardIDsFromEveryNode(t *testing.T) {
	runContainers(2, 4)
	for port := range containersInfos {
		res := sendRequest(port, "GET", "shard/all_ids", t, http.StatusOK, nil)
		_, ok := res["shard_ids"]
		if !ok {
			t.Errorf("localhost:%v: empty shard_ids received", port)
			t.FailNow()
		}
		var expectedIDs string
		for i := 0; i < initNumOfShards; i++ {
			expectedIDs += strconv.Itoa(i) + ","
		}
		expectedIDs = strings.TrimRight(expectedIDs, ",")

		if expectedIDs != res["shard_ids"].(string) {
			t.Errorf("localhost:%v: Expected: %v: %v", port, expectedIDs, res["shard_ids"])
		}
		if res["result"].(string) != "Success" {
			t.Errorf("localhost:%v: result must be Success: %v", port, res["result"].(string))
		}
	}
}

func TestMinNumMembersInShard(t *testing.T) {
	runContainers(2, 4)
	t.Run("Min", chechMinNumMembersInShardFromEveryNode)
}

func TestMinNumMembersInShardAfterAddingOneNode(t *testing.T) {
	runContainers(2, 4)
	newNetworkIP := runAContainer()
	putView := map[string]interface{}{"result": "Success", "msg": "Successfully added " + newNetworkIP + ":8080 to view"}
	checkView(t, "PUT", putView, newNetworkIP, !askEveryNode)
	chechMinNumMembersInShardFromEveryNode(t)
}

func TestKeyNotExist(t *testing.T) {
	runContainers(2, 4)
	expectedGetKey := map[string]interface{}{"result": "Error", "msg": "Key does not exist", "statusCode": http.StatusNotFound}
	checkGetKey(t, "bumba", expectedGetKey, askEveryNode)
}

func TestPutKeyFalse(t *testing.T) {
	runContainers(2, 4)
	expectedPutKey := map[string]interface{}{"replaced": false, "msg": "Added successfully", "statusCode": http.StatusCreated}
	checkPutKey(t, "Kuku", "Monkey", expectedPutKey, !askEveryNode)
}

func Test5Shard4Nodes(t *testing.T) {
	runContainers(1, 4)
	chechMinNumMembersInShardFromEveryNode(t)
	putShard := map[string]interface{}{"result": "Error", "msg": "Not enough nodes for 5 shards", "statusCode": http.StatusBadRequest}
	checkPutShard(t, 5, putShard)
}

func Test3Shard4Nodes(t *testing.T) {
	runContainers(1, 4)
	chechMinNumMembersInShardFromEveryNode(t)
	putShard := map[string]interface{}{"result": "Error", "msg": "Not enough nodes. 3 shards result in a nonfault tolerant shard", "statusCode": http.StatusBadRequest}
	checkPutShard(t, 3, putShard)
}

func Test2Shard4Nodes(t *testing.T) {
	runContainers(1, 4)
	chechMinNumMembersInShardFromEveryNode(t)
	putShard := map[string]interface{}{"result": "Success", "statusCode": http.StatusOK}
	checkPutShard(t, 2, putShard)
	chechMinNumMembersInShardFromEveryNode(t)
}
func Test1Shard4Nodes(t *testing.T) {
	runContainers(2, 4)
	chechMinNumMembersInShardFromEveryNode(t)
	putShard := map[string]interface{}{"result": "Success", "statusCode": http.StatusOK}
	checkPutShard(t, 1, putShard)
	chechMinNumMembersInShardFromEveryNode(t)
}

//--------------------------------------- Helper Functions----------------------------------/

func initPutKey(t *testing.T, key string, value string) {
	expectedPutKey := map[string]interface{}{"replaced": false, "msg": "Added successfully", "statusCode": http.StatusCreated}
	checkPutKey(t, key, value, expectedPutKey, !askEveryNode)
}

func deleteNode(t *testing.T, node string, deleteErrbody bool) {
	body := url.Values{}
	body.Add("ip_port", node)
	for port := range containersInfos {
		res := sendRequest(port, "DELETE", "view", t, http.StatusOK, body)
		log.Println("DELETE RES: ", res)
		if !deleteErrbody {
			return
		}
	}
	t.Logf("Before deleting node: %v", node)
	delete(containersInfos, node)
	t.Logf("After deleting node: %v", node)

}
func chechMinNumMembersInShardFromEveryNode(t *testing.T) {
	minNumMembers := int(math.Ceil(float64(initNumOfContainers) / float64(initNumOfShards)))
	for port := range containersInfos {
		for shardID := 0; shardID < initNumOfShards; shardID++ {
			res := sendRequest(port, "GET", "shard/members/"+strconv.Itoa(shardID), t, http.StatusOK, nil)
			_, ok := res["members"]
			if !ok {
				t.Errorf("localhost:%v: members must not be empty for shard_id %v: %v", port, shardID, res)
				t.FailNow()
			} else {
				members := strings.Split(res["members"].(string), ",")
				if len(members) < minNumMembers {
					t.Errorf("localhost:%v: not enough members in shard_id: %v. Members: %v", port, shardID, members)
				}
			}
		}
	}
}

func checkView(t *testing.T, typeReq string, expected map[string]interface{}, newNetworkIP string, askEveryone bool) {
	for port := range containersInfos {
		body := url.Values{}
		body.Add("ip_port", newNetworkIP+":8080")
		res := sendRequest(port, typeReq, "view", t, http.StatusOK, body)
		_, ok := res["msg"].(string)
		if !ok {
			t.Errorf("localhost:%v: unable to put view network ip:port: %v:8080", port, newNetworkIP)
			t.FailNow()
		} else {
			if res["msg"] != expected["msg"] {
				t.Errorf("localhost:%v: Expected: %v: %v", port, expected["msg"], res["msg"])
				t.FailNow()
			}
			if res["result"] != expected["result"] {
				t.Errorf("localhost:%v: Expected: %v : %v", port, expected["result"], res["result"])
				t.FailNow()
			}
		}
		if !askEveryone {
			break
		}
	}
}

func checkGetKey(t *testing.T, key string, expected map[string]interface{}, askEveryone bool) float64 {
	var ownerID = -1.
	for port := range containersInfos {
		body := url.Values{}
		res := sendRequest(port, "GET", "keyValue-store/"+key, t, expected["statusCode"].(int), body)
		// log.Printf("Port: %v. GET RES: %v\n", port, res)
		_, ok := res["result"].(string)
		if !ok {
			t.Errorf("localhost:%v: unable to get key: %v", port, key)
			t.FailNow()
		} else {
			if res["result"] != expected["result"] {
				t.Errorf("localhost:%v: Expected: %v : %v", port, expected["result"], res["result"])
				t.FailNow()
			}
			if expected["result"] == "Success" {
				if res["msg"] != expected["msg"] {
					t.Errorf("localhost:%v: Expected: %v: %v", port, expected["msg"], res["msg"])
					t.FailNow()
				}
				owner, isExist := res["owner"]
				if !isExist {
					t.Errorf("localhost:%v owner is not exist", port)
					t.FailNow()
				}
				ownerID = owner.(float64)

			} else if res["msg"] != expected["msg"] {
				t.Errorf("localhost:%v: Expected %v: %v", port, expected["msg"], res["msg"])
			}
		}
		if !askEveryone {
			break
		}
	}
	return ownerID
}

func checkDeleteKey(t *testing.T, key string, expected map[string]interface{}, deleteFromEveryone bool) {
	for port := range containersInfos {
		body := url.Values{}
		res := sendRequest(port, "DELETE", "keyValue-store/"+key, t, expected["statusCode"].(int), body)
		msg, _ := res["msg"].(string)
		if msg != expected["msg"].(string) {
			t.Errorf("localhost:%v: Expected: %v: %v", port, expected["msg"], msg)
		}
		if !deleteFromEveryone {
			return
		}
	}
}

func checkPutKey(t *testing.T, key string, val string, expected map[string]interface{}, putEveryone bool) {
	for port := range containersInfos {
		body := url.Values{}
		body.Add("val", val)
		res := sendRequest(port, "PUT", "keyValue-store/"+key, t, expected["statusCode"].(int), body)
		t.Log(res)
		replaced, isExistReplaced := res["replaced"].(bool)
		msg, isExistMsg := res["msg"].(string)
		if !isExistMsg || !isExistReplaced {
			t.Errorf("localhost:%v: msg and replace not exist", port)
			t.FailNow()
		} else {
			if replaced != expected["replaced"].(bool) {
				t.Errorf("localhost:%v: Expected: %v: %v", port, expected["replaced"], replaced)
			}
			if msg != expected["msg"].(string) {
				t.Errorf("localhost:%v: Expected: %v: %v", port, expected["msg"], msg)
			}
		}
		if !putEveryone {
			return
		}
	}
}

func getCount(t *testing.T, shardID float64, statusCode int, checkEveryone bool) float64 {
	var count float64

	shardIDString := fmt.Sprint(shardID)
	for port := range containersInfos {
		body := url.Values{}
		res := sendRequest(port, "GET", "shard/count/"+shardIDString, t, statusCode, body)
		localCount, _ := res["Count"].(float64)
		if localCount >= count {
			count = localCount
		}
		if !checkEveryone {
			break
		}
	}

	return count
}

func putKey(t *testing.T, port string, key string, val string, statusCode int) map[string]interface{} {
	body := url.Values{}
	body.Add("val", val)
	res := sendRequest(port, "PUT", "keyValue-store/"+key, t, statusCode, body)
	return res
}

func checkPutShard(t *testing.T, num int, expected map[string]interface{}) {
	for port := range containersInfos {
		body := url.Values{}
		body.Add("num", strconv.Itoa(num))
		res := sendRequest(port, "PUT", "shard/changeShardNumber", t, expected["statusCode"].(int), body)
		result, isExist := res["result"]
		if !isExist {
			t.Errorf("localhost:%v: received empty body: %v", port, res)
		} else {
			if result == "Success" {
				initNumOfShards = num
				var expectedIDs string
				for i := 0; i < initNumOfShards; i++ {
					expectedIDs += strconv.Itoa(i) + ","
				}
				expectedIDs = strings.TrimRight(expectedIDs, ",")
				if expectedIDs != res["shard_ids"].(string) {
					t.Errorf("localhost:%v: Expected: %v: %v", port, expectedIDs, res["shard_ids"])
				}
				return
			}
			if result != expected["result"] {
				t.Errorf("localhost:%v: Expected result: %v: %v", port, expected["result"], result)
			}
			if res["msg"] != expected["msg"] {
				t.Errorf("localhost:%v: Expected msg: %v: %v", port, expected["msg"], res["msg"])
			}
		}
		return
	}
}
