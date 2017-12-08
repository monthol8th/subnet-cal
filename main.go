package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type PossibleIPStruct struct {
	NetworkAddress   string
	BroadcastAddress string
	Usable           string
}

type ResponsePayload struct {
	Status           string
	IP               string
	Subnet           string
	BSubnet          string
	NetworkAddress   string
	BroadcastAddress string
	NumberOfHost     uint64
	Usable           string
	BinID            string
	IntID            uint64
	HexID            string
	Possible         []PossibleIPStruct
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

func ipArrayToString(ip [4]uint8) string {
	var buffer bytes.Buffer
	for i, v := range ip {
		if i > 0 {
			buffer.WriteString(".")
		}
		buffer.WriteString(strconv.Itoa(int(v)))
	}
	return buffer.String()
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

func firstlastHost(network [4]uint8, broadcast [4]uint8) string {
	first := network
	last := broadcast
	first[3]++
	last[3]--
	fStr := ipArrayToString(first)
	lStr := ipArrayToString(last)
	var buffer bytes.Buffer
	buffer.WriteString(fStr)
	buffer.WriteString(" - ")
	buffer.WriteString(lStr)
	return buffer.String()
}

func numberOfHostCalculate(subnet [4]uint8) (number uint64) {
	number = 0
	for _, value := range subnet {
		number *= 256
		number += uint64(^value)
	}
	number++
	return
}

func calc(ip [4]uint8, subnet [4]uint8) (network, broadcast [4]uint8) {
	network = networkAddress(ip, subnet)
	broadcast = broadcastAddress(ip, subnet)
	return
}

func possibleRange(ip [4]uint8, subnet [4]uint8, numberOfHost uint64) (everyRange []PossibleIPStruct) {
	var addr uint8
	var inc uint16
	possibleIP := ip
	isAddrFound := false
	for i, v := range subnet {
		if isAddrFound {
			possibleIP[i] = 0
		} else {
			bwSubnet := ^v
			if bwSubnet != 0 {
				addr = uint8(i)
				inc = uint16(bwSubnet) + 1
				isAddrFound = true
				possibleIP[i] = 0
			}
		}
	}
	numberOfPossible := 256 / uint16(inc)
	for i := uint16(0); i < numberOfPossible; i++ {
		possibleNetwork, possibleBroadcast := calc(possibleIP, subnet)
		possibleIP[addr] += uint8(inc)
		usable := "none"
		if numberOfHost-2 > 0 {
			usable = firstlastHost(possibleNetwork, possibleBroadcast)
		}
		everyRange = append(everyRange, PossibleIPStruct{
			NetworkAddress:   ipArrayToString(possibleNetwork),
			BroadcastAddress: ipArrayToString(possibleBroadcast),
			Usable:           usable})
	}
	return
}

func bit8(x string) string {
	l := 8 - len(x)
	var b bytes.Buffer
	b.WriteString(strings.Repeat("0", l))
	b.WriteString(x)
	return b.String()
}

func ipToBinary(ip [4]uint8) string {
	var b bytes.Buffer
	for i, v := range ip {
		bit8ip := bit8(strconv.FormatInt(int64(v), 2))
		if i > 0 {
			b.WriteString(".")
			b.WriteString(bit8ip)
		}
	}
	return b.String()
}

func toBinary(ip [4]uint8) string {
	var b bytes.Buffer
	for _, v := range ip {
		b.WriteString(bit8(strconv.FormatInt(int64(v), 2)))
	}
	return b.String()
}
func root(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	ipString := r.PostFormValue("ip")
	subnetString := r.PostFormValue("subnet")

	ipArray := ipStringToUint(ipString)
	subnetArray := ipStringToUint(subnetString)
	networkAddressArray, broadcastAddressArray := calc(ipArray, subnetArray)
	numberOfHost := numberOfHostCalculate(subnetArray)
	possible := possibleRange(ipArray, subnetArray, numberOfHost)

	var res ResponsePayload
	res.Status = "OK"
	res.IP = ipString
	res.Subnet = subnetString
	res.NetworkAddress = ipArrayToString(networkAddressArray)
	res.BroadcastAddress = ipArrayToString(broadcastAddressArray)
	res.NumberOfHost = numberOfHost
	res.BSubnet = ipToBinary(subnetArray)
	res.BinID = toBinary(ipArray)
	res.IntID, _ = strconv.ParseUint(res.BinID, 2, 64)
	res.HexID = strconv.FormatUint(res.IntID, 16)
	if res.Usable = "none"; numberOfHost-2 > 0 {
		res.Usable = firstlastHost(networkAddressArray, broadcastAddressArray)
	}
	res.Possible = possible
	jsonRes, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(jsonRes))
}

func main() {
	http.HandleFunc("/", root)
	fmt.Println("listening and serving on PORT 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
