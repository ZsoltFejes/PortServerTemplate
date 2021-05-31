package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Status    string
	Completed bool
	ID        string
}

// Start REST Api Server
func startHttpServer() {
	l("Starting API server...", false, true)
	port := ":8080"
	http.HandleFunc("/", testBroadcast)
	l("API server is listening on port "+port, false, true)
	log.Fatal(http.ListenAndServe(port, nil))
}

func testBroadcast(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		id := getID()
		server.broadcast <- Job{Message: "Hello from API", ID: id}
		w.WriteHeader(http.StatusOK)
		response := Response{Status: "Broadcasted", Completed: true, ID: id}
		json.NewEncoder(w).Encode(response)
	default:
		w.WriteHeader(http.StatusNotFound)
		response := Response{Status: "Error", Completed: false}
		json.NewEncoder(w).Encode(response)
	}
}
