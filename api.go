package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Response struct {
	Status    string
	Completed bool
}

// Start REST Api Server
func startHttpServer() {
	http.HandleFunc("/", testBroadcast)
	http.HandleFunc("/broadcast", testBroadcast)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// TODO: This API end point have not been tested
func testBroadcast(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		broadcastMessage := Command{}
		err := json.NewDecoder(r.Body).Decode(&broadcastMessage)
		if err != nil {
			fmt.Println("Unable to parse Request")
			w.WriteHeader(http.StatusBadRequest)
			response := Response{Status: "Error", Completed: false}
			json.NewEncoder(w).Encode(response)
		} else {
			if broadcastMessage.Broadcast != "" {
				fmt.Println("Broadcast message has been received")
				manager.broadcast <- broadcastMessage
				w.WriteHeader(http.StatusOK)
				response := Response{Status: "Broadcasted", Completed: true}
				json.NewEncoder(w).Encode(response)
			} else {
				fmt.Println("Broadcast Field was not found")
				w.WriteHeader(http.StatusBadRequest)
				response := Response{Status: "Error", Completed: false}
				json.NewEncoder(w).Encode(response)
			}

		}
	}
}
