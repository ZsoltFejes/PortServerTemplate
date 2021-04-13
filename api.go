package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Response struct {
	Status    string
	Completed bool
}

// Start REST Api Server
func startHttpServer() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", testBroadcast)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func testBroadcast(w http.ResponseWriter, r *http.Request) {
	testCommand := Command{Command: "Hello From API"}
	manager.broadcast <- testCommand
	w.WriteHeader(http.StatusOK)
	response := Response{Status: "Broadcasted", Completed: true}
	json.NewEncoder(w).Encode(response)
}
