package main

import (
	"encoding/json"
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
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func testBroadcast(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		testCommand := Command{Command: "Hello From API"}
		manager.broadcast <- testCommand
		w.WriteHeader(http.StatusOK)
		response := Response{Status: "Broadcasted", Completed: true}
		json.NewEncoder(w).Encode(response)
	}
}
