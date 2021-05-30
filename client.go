package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

type Client struct {
	socket net.Conn
	data   chan Job
}

// Handle received messages in clinet mode
func (client *Client) receive() {
	var job Job
	decoder := json.NewDecoder(client.socket)
	for {
		err := decoder.Decode(&job)
		if err != nil {
			client.socket.Close()
			break
		}
		handleJob(&job, client)
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
		input, _ := reader.ReadString('\n')
		args := strings.Fields(input)
		if len(args) > 0 {
			command := Job{Command: args[0], ID: getID()}
			err := encoder.Encode(command)
			if err != nil {
				fmt.Printf("Encoding Error: %s", err)
				client.socket.Close()
				break
			}
		} else {
			fmt.Println("Type a command")
		}
	}
}
