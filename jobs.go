package main

import (
	"fmt"
)

type Job struct {
	Command   string `json:"command,omitempty"`
	Message   string `json:"message,omitempty"`
	Status    string `json:"status,omitempty"`
	Broadcast string `json:"broadcast,omitempty"`
}

// Handle incoming jobs
func handleJob(job *Job, client *Client) {
	if len(job.Command) > 0 {
		fmt.Printf("Received Job: %s\n", job.Command)
		// Handle different commands by calling a function
		switch job.Command {
		case "ping":
			pong(client)
		}
	}
	if len(job.Status) > 0 {
		fmt.Printf("Status: %s\n", job.Status)
	}
	if len(job.Broadcast) > 0 {
		fmt.Println("Broadcast Message: " + job.Broadcast)
	}
	if len(job.Message) > 0 {
		fmt.Println("Message: " + job.Message)
	}
}

func pong(client *Client) {
	pong := Job{Message: "PONG"}
	client.data <- pong
}
