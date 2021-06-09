package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"log"
	"net"
	"os"
	"strings"
	"time"
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
		client.comandSent <- false
		go handleJob(job, client)
		job.reset()
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
	client := &Client{comandSent: make(chan bool)}
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
	go client.latency() // Latency is added for testing
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
			client.comandSent <- true
		} else {
			l("Type a command", false, true)
		}
	}
}

// This feature will be part of job processor instead of the client
func (client *Client) latency() {
	var now time.Time
	waiting := false
	for {
		select {
		case commandSent := <-client.comandSent:
			if commandSent {
				now = time.Now()
				if waiting {
					l("Application was waiting for an answer from the server", false, true)
				}
				waiting = true
			} else {
				waiting = false
				l(time.Since(now).String(), false, true)
			}
		}
	}
}
