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
	Status string
	Ip     [4]uint8
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

func root(w http.ResponseWriter, r *http.Request) {
	ipString := r.PostFormValue("ip")
	ipArray := ipStringToUint(ipString)

	var res ResponsePayload
	res.Status = "OK"
	res.Ip = ipArray
	jsonRes, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(jsonRes))
}

func main() {
	http.HandleFunc("/", root)
	fmt.Println("listening and serving on PORT 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
