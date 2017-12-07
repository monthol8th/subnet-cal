package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ResponsePayload struct {
	Status string
}

func homePage(w http.ResponseWriter, r *http.Request) {
	var m ResponsePayload
	m.Status = "OK"
	b, _ := json.Marshal(m)

	w.Header().Set("Content-Type", "application/json")

	fmt.Fprintf(w, string(b))
}

func main() {
	http.HandleFunc("/", homePage)
	fmt.Println("listening and serving on PORT 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
