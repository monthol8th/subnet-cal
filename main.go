package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type ResponsePayload struct {
	Status           string
	IP               [4]uint8
	Subnet           [4]uint8
	NetworkAddress   [4]uint8
	BroadcastAddress [4]uint8
}

func homePage(w http.ResponseWriter, r *http.Request) {
	var m ResponsePayload
	m.Status = "OK"
	b, _ := json.Marshal(m)

	w.Header().Set("Content-Type", "application/json")

	fmt.Fprintf(w, string(b))
}

func ipStringToUint(ipString string) (ip [4]uint8) {
	splitedIPString := strings.Split(ipString, ".")
	for i, v := range splitedIPString {
		temp, _ := strconv.ParseInt(v, 10, 64)
		ip[i] = uint8(temp)
	}
	return
}

func networkAddress(ip [4]uint8, subnet [4]uint8) (addr [4]uint8) {
	for i := range ip {
		addr[i] = ip[i] & subnet[i]
	}
	return
}

func broadcastAddress(ip [4]uint8, subnet [4]uint8) (addr [4]uint8) {
	for i := range ip {
		addr[i] = ip[i] | ^subnet[i]
	}
	return
}

func root(w http.ResponseWriter, r *http.Request) {
	ipString := r.PostFormValue("ip")
	subnetString := r.PostFormValue("subnet")

	ipArray := ipStringToUint(ipString)
	subnetArray := ipStringToUint(subnetString)
	networkAddressArray := networkAddress(ipArray, subnetArray)
	broadcastAddressArray := broadcastAddress(ipArray, subnetArray)

	var res ResponsePayload
	res.Status = "OK"
	res.IP = ipArray
	res.Subnet = subnetArray
	res.NetworkAddress = networkAddressArray
	res.BroadcastAddress = broadcastAddressArray
	jsonRes, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(jsonRes))
}

func main() {
	http.HandleFunc("/", root)
	fmt.Println("listening and serving on PORT 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
