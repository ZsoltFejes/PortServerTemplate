package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type Client struct {
	socket net.Conn
	data   chan Command
}

// Handle received messages in clinet mode
func (client *Client) receive() {
	var command Command
	decoder := json.NewDecoder(client.socket)
	for {
		err := decoder.Decode(&command)
		if err != nil {
			client.socket.Close()
			break
		}
		handleCommand(&command, client)
	}
}

func startClientMode(encrypt *bool) {
	fmt.Println("Starting client...")
	client := &Client{}
	if *encrypt {
		// For Testing certificate verification is disabled
		config := tls.Config{InsecureSkipVerify: true}
		connection, err := tls.Dial("tcp", "localhost:12345", &config)
		checkErr("Connecting to server with TLS error", err)
		client.socket = connection
	} else {
		connection, err := net.Dial("tcp", "localhost:12345")
		checkErr("Connecting to server error", err)
		client.socket = connection
	}
	go client.receive()
	encoder := json.NewEncoder(client.socket)
	for {
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')
		fmt.Printf("%s", message)
		// Test Strings:
		testMessage := Command{Command: "Hello"}
		err := encoder.Encode(testMessage)
		if err != nil {
			fmt.Printf("Encoding Error: %s", err)
			client.socket.Close()
			break
		}
	}
}
