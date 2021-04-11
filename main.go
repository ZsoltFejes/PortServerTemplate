package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan Command
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	socket net.Conn
	data   chan Command
}

type Command struct {
	Command string `json:"command,omitempty"`
	Status  string `json:"status,omitempty"`
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

// Handle receaved messages in server, needs to be assnged to each client, curently expecting only json strings
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
		handleCommand(&command)
		// manager.broadcast <- message
		acknowledge := Command{Status: "Acknowledged"}
		client.data <- acknowledge
	}
}

// Handle receaved messages in clinet mode
func (client *Client) receive() {
	var command Command
	decoder := json.NewDecoder(client.socket)
	for {
		err := decoder.Decode(&command)
		if err != nil {
			client.socket.Close()
			break
		}
		handleCommand(&command)
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

func handleCommand(command *Command) {
	if len(command.Command) > 0 {
		fmt.Printf("Command: %s\n", command.Command)
	}
	if len(command.Status) > 0 {
		fmt.Printf("Status: %s\n", command.Status)
	}
}

func startServerMode() {
	fmt.Println("Starting server...")
	listener, err := net.Listen("tcp", ":12345")
	if err != nil {
		fmt.Println(err)
	}
	manager := ClientManager{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan Command),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	go manager.start()
	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}
		client := &Client{socket: connection, data: make(chan Command)}
		manager.register <- client
		go manager.receive(client)
		go manager.send(client)
	}
}

func startClientMode() {
	fmt.Println("Starting client...")
	connection, err := net.Dial("tcp", "localhost:12345")
	if err != nil {
		fmt.Println(err)
	}
	client := &Client{socket: connection}
	go client.receive()
	encoder := json.NewEncoder(client.socket)
	for {
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')
		fmt.Printf("%s", message)
		// Test Strings:
		testMessage := Command{Command: "Hello"}
		encoder.Encode(testMessage)
	}
}

func main() {
	flagMode := flag.String("mode", "server", "Start in client or server mode")
	flag.Parse()
	if strings.ToLower(*flagMode) == "server" {
		startServerMode()
	} else {
		startClientMode()
	}
}
