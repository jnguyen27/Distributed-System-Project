/**********************************************
/ AUTHOR: Akobir Khamidov (akhamido@ucsc.edu) *
/ AUTHOR: Vien Van (vhvan@ucsc.edu)			  *
/ COPYRIGHT 2018 © by TEAMAWESOME			  *
***********************************************/

package hw4

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

// ------------ CHANGE THIS LINE --------------------------\\
var subnetName = "mynet"
var prefixSubnetAdress = "10.0.0."
var prefixPort = "80"
var path = "../src/." // tell where your docker file EX: "." Current path. "../." One path down
// --------------------END ----------------------------------\\

// ----------------------------------------------------- DO NOT GO BEYOND THIS LINE ----------------------------------------------------\\

//------------------ - Global Variables----------------------\\
var nextID = 10
var containersInfos map[string]map[string]string
var askEveryNode = true
var initNumOfShards int
var initNumOfContainers int // Error if > 8

//----------------------END-----------------------------------\\
func init() {
	buildFlag := flag.Bool("build", false, "a bool")
	flag.Parse()
	removeAllContainers()
	if *buildFlag {
		exec.Command("docker", "rmi", "-f", "testing").Run()
		cmd := exec.Command("docker", "build", "-t", "testing", path)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			log.Fatal("Docker FAILED: Make sure docker is running or no errors from programm")
		} else {
			log.Println("Successfully created image 'Testing'")
		}
	}

}
func runContainers(numOfShards int, numOfContainers int) {
	killContainers()
	initNumOfShards = numOfShards
	initNumOfContainers = numOfContainers
	initContainersInfos()
	for port := range containersInfos {
		runAContainer(port)
	}
}

func initContainersInfos() {
	containersInfos = make(map[string]map[string]string)
	for i := 10; i < 10+initNumOfContainers; i++ {
		key := prefixPort + strconv.Itoa(i)
		containersInfos[key] = make(map[string]string)
		containersInfos[key]["networkIp"] = prefixSubnetAdress + strconv.Itoa(i)
	}
}

func generateView() string {
	var str string
	for _, container := range containersInfos {
		str += container["networkIp"] + ":8080" + ","
	}
	str = strings.TrimRight(str, ",")
	return str
}

func removeAllContainers() {
	cmd := exec.Command("docker", "ps", "-aq")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal("Docker FAILED: Make sure docker is running or no errors from programm")
	}
	containerIDs := strings.Split(out.String(), "\n")
	for _, containerID := range containerIDs {
		exec.Command("docker", "rm", "-f", containerID).Run()
	}
	log.Println("Removed all containers")
}

func killContainers() {
	cmd := exec.Command("docker", "ps", "-aq")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()
	containerIDs := strings.Split(out.String(), "\n")
	for _, containerID := range containerIDs {
		exec.Command("docker", "kill", containerID).Run()
	}
}

func runAContainer(ports ...string) string {
	var port string
	if len(ports) < 1 {
		myID := strconv.Itoa(getNextID())
		port = prefixPort + myID
		containersInfos[port] = make(map[string]string)
		containersInfos[port]["networkIp"] = prefixSubnetAdress + myID
	} else {
		port = ports[0]
	}
	container := containersInfos[port]
	args := []string{"run", "-p", port + ":8080", "--net=" + subnetName, "--ip=" + container["networkIp"], "-e", "VIEW=" + generateView(),
		"-e", "IP_PORT=" + container["networkIp"] + ":8080", "-e", "S=" + strconv.Itoa(initNumOfShards), "-d", "testing"}
	cmd := exec.Command("docker", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal("Unable to run container. Make sure Docker is running", err, cmd.Args)
	}
	containersInfos[port]["id"] = out.String()

	nextID++

	return containersInfos[port]["networkIp"]
}

func sendRequest(port string, typeReq string, route string, t *testing.T, statusCode int, bodyReq url.Values) map[string]interface{} {
	client := http.Client{}
	url := "http://localhost:" + port + "/" + route
	request, err := http.NewRequest(typeReq, url, strings.NewReader(bodyReq.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatalf("Unable to connect to %v. TEST TERMINATED", port)
	}
	resp, err := client.Do(request)
	if err != nil {
		t.Fatalf("Unable to receive response from %v.  TEST TERMINATED", port)
	}
	var body map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&body)
	defer resp.Body.Close()
	if resp.StatusCode != statusCode {
		t.Errorf("Status code not match from localhost:%v. Received %v Expected %v", port, resp.StatusCode, statusCode)
	}
	return body
}

func getNextID() int {
	return nextID
}

func getPort() string {
	var node string
	for port := range containersInfos {
		networkIP := containersInfos[port]["networkIp"]
		node = networkIP + ":" + "8080"
		break
	}
	return node
}
