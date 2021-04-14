package main

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
)

type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan Command
	register   chan *Client
	unregister chan *Client
}

var manager = ClientManager{
	clients:    make(map[*Client]bool),
	broadcast:  make(chan Command),
	register:   make(chan *Client),
	unregister: make(chan *Client),
}

// Start Client manager
func (manager *ClientManager) start() {
	for {
		select {
		case connection := <-manager.register:
			manager.clients[connection] = true
			fmt.Println("Added new connection!")
		case connection := <-manager.unregister:
			if _, ok := manager.clients[connection]; ok {
				close(connection.data)
				delete(manager.clients, connection)
				fmt.Println("A connection has been terminated")
			}
		case message := <-manager.broadcast:
			for connection := range manager.clients {
				select {
				case connection.data <- message:
				default:
					close(connection.data)
					delete(manager.clients, connection)
				}
			}
		}
	}
}

// Handle received messages in server, needs to be assinged to each client,
// curently expecting only json strings
func (manager *ClientManager) receive(client *Client) {
	var command Command
	decoder := json.NewDecoder(client.socket)
	for {
		err := decoder.Decode(&command)
		if err != nil {
			fmt.Println(err)
			manager.unregister <- client
			client.socket.Close()
			break
		}
		fmt.Println("Received commands")
		handleCommand(&command, client)
	}
}

// Send message to client
func (manager *ClientManager) send(client *Client) {
	defer client.socket.Close()
	encoder := json.NewEncoder(client.socket)
	for {
		select {
		case command, ok := <-client.data:
			if !ok {
				return
			}
			encoder.Encode(command)
		}
	}
}

func startServerMode(manager *ClientManager, ecrypt *bool) {
	fmt.Println("Starting server...")
	var listener net.Listener
	if *ecrypt {
		basePath := os.Getenv("GOPATH") + "/src/github.com/ZsoltFejes/go_link/"
		cert, err := tls.LoadX509KeyPair(basePath+"server.crt", basePath+"server.key")
		checkErr("Importing TLS certifiacets", err)
		config := tls.Config{Certificates: []tls.Certificate{cert}}
		now := time.Now()
		config.Time = func() time.Time { return now }
		config.Rand = rand.Reader
		listener, err = tls.Listen("tcp", ":12345", &config)
		checkErr("Creating TLS listener", err)
	} else {
		var err error
		listener, err = net.Listen("tcp", ":12345")
		checkErr("Creating NET listener", err)
	}
	go manager.start()
	go startHttpServer()
	for {
		connection, err := listener.Accept()
		checkErr("Accepting connection", err)
		client := &Client{socket: connection, data: make(chan Command)}
		manager.register <- client
		go manager.receive(client)
		go manager.send(client)
	}
}
