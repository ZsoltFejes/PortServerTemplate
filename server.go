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

type Server struct {
	clients    map[*Client]bool
	broadcast  chan Job
	register   chan *Client
	unregister chan *Client
}

var server = Server{
	clients:    make(map[*Client]bool),
	broadcast:  make(chan Job),
	register:   make(chan *Client),
	unregister: make(chan *Client),
}

// Start Client manager
func (server *Server) start() {
	for {
		select {
		case connection := <-server.register:
			server.clients[connection] = true
			fmt.Println("Added new connection! " + connection.socket.RemoteAddr().String())
		case connection := <-server.unregister:
			if _, ok := server.clients[connection]; ok {
				close(connection.data)
				delete(server.clients, connection)
				fmt.Println("A connection has been terminated " + connection.socket.RemoteAddr().String())
			}
		case message := <-server.broadcast:
			for connection := range server.clients {
				select {
				case connection.data <- message:
				default:
					close(connection.data)
					delete(server.clients, connection)
				}
			}
		}
	}
}

// Handle received messages in server, needs to be assinged to each client,
// curently expecting only json strings
func (server *Server) receive(client *Client) {
	var job Job
	decoder := json.NewDecoder(client.socket)
	for {
		err := decoder.Decode(&job)
		if err != nil {
			fmt.Println(err)
			server.unregister <- client
			client.socket.Close()
			break
		}
		handleJob(&job, client)
	}
}

// Send message to client
func (server *Server) send(client *Client) {
	defer client.socket.Close()
	encoder := json.NewEncoder(client.socket)
	for {
		select {
		case job, ok := <-client.data:
			if !ok {
				return
			}
			encoder.Encode(job)
		}
	}
}

func startServerMode(manager *Server, ecrypt *bool) {
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
		client := &Client{socket: connection, data: make(chan Job)}
		manager.register <- client
		go manager.receive(client)
		go manager.send(client)
	}
}
