package main

import (
	"flag"
	"fmt"
	"strings"
)

type Job struct {
	Command   string `json:"Job,omitempty"`
	Message   string `json:"Message,omitempty"`
	Status    string `json:"status,omitempty"`
	Broadcast string `json:"broadcast,omitempty"`
}

func handleJob(job *Job, client *Client) {
	if len(job.Command) > 0 {
		fmt.Printf("Received Job: %s\n", job.Command)
		acknowledge := Job{Status: "Acknowledged"}
		client.data <- acknowledge
	}
	if len(job.Status) > 0 {
		fmt.Printf("Status: %s\n", job.Status)
	}
	if len(job.Broadcast) > 0 {
		fmt.Println("Broadcast Message: " + job.Broadcast)
	}
}

// Check Error function for universal error handling
func checkErr(message string, err interface{}) {
	if err != nil {
		fmt.Printf("%s- %s\n", message, err)
	}
}

func main() {
	flagMode := flag.String("mode", "server", "Start in client or server mode")
	flagTLS := flag.Bool("tls", false, "Set Server to use TLS (Add certifiacet to root directory as server.crt and server.key)")
	flag.Parse()
	if strings.ToLower(*flagMode) == "server" {
		startServerMode(&manager, flagTLS)
	} else {
		startClientMode(flagTLS)
	}
}
