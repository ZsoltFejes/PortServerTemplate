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
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	socket net.Conn
	data   chan []byte
}

type Command struct {
	Command string `json:"command"`
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
	for {
		var commands map[string]Command
		decoder := json.NewDecoder(client.socket)
		err := decoder.Decode(&commands)
		if err != nil {
			manager.unregister <- client
			client.socket.Close()
			break
		}
		fmt.Println("Received commands")
		handleCommands(&commands)
		// manager.broadcast <- message
		client.data <- []byte(`{"status":"acknowledged"}`)
	}
}

// Handle receaved messages in clinet mode
func (client *Client) receive() {
	for {
		message := make([]byte, 4096)
		length, err := client.socket.Read(message)
		if err != nil {
			client.socket.Close()
			break
		}
		if length > 0 {
			fmt.Println("RECEIVED: " + string(message))
		}
	}
}

// Send message to client
func (manager *ClientManager) send(client *Client) {
	defer client.socket.Close()
	for {
		select {
		case message, ok := <-client.data:
			if !ok {
				return
			}
			client.socket.Write(message)
		}
	}
}

func handleCommands(commands *map[string]Command) {
	for cKey, command := range *commands {
		fmt.Printf("Key: %s, Command: %s\n", cKey, command)
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
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	go manager.start()
	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}
		client := &Client{socket: connection, data: make(chan []byte)}
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
	for {
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')
		fmt.Println(message)
		// Test Strings:
		testMessage := `{ "1": {"command": "Hello"}, "2": {"command": "Hello2"}}`
		connection.Write([]byte(strings.TrimRight(testMessage, "\n")))
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
