package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
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
		go handleJob(job, client)
		job.reset()
	}
}

func startClientMode(encrypt *bool) {
	l("Starting client...", false, true)
	client := &Client{}
	if *encrypt {
		// For Testing certificate verification is disabled
		config := tls.Config{InsecureSkipVerify: true}
		connection, err := tls.Dial("tcp", appConfig.Api.Address+":"+appConfig.Client.Port, &config)
		checkErr("Connecting to server with TLS error", err)
		client.socket = connection
	} else {
		connection, err := net.Dial("tcp", appConfig.Api.Address+":"+appConfig.Client.Port)
		checkErr("Connecting to server error", err)
		client.socket = connection
	}
	l("Client Connected", false, true)
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
				l("Encoding Error: "+err.Error(), false, false)
				client.socket.Close()
				break
			}
		} else {
			l("Type a command", false, true)
		}
	}
}
