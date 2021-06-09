package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"log"
	"net"
	"os"
	"strings"
)

// TODO Rework latency and outgoing requests to be handeled per job, so outgoung requests will be part of the jpb processor

type Client struct {
	socket     net.Conn
	data       chan Job
	comandSent chan bool
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

func (client *Client) send() {
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

func startClientMode() {
	// Create or open log directory
	f, err := os.OpenFile(WORKDIR+`/client.log`, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		l(err.Error(), true, true)
	}
	defer f.Close()

	log.SetOutput(f)
	l("Starting client...", false, true)
	client := &Client{
		comandSent: make(chan bool),
		data:       make(chan Job)}
	if appConfig.Tls {
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
	go client.send()
	for {
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		args := strings.Fields(input)
		if len(args) > 0 {
			job := Job{Command: args[0], Args: args[1:], Client: client}
			job.new()
		} else {
			l("Type a command", false, true)
		}
	}
}
