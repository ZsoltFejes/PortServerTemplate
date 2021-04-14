package main

import (
	"flag"
	"fmt"
	"strings"
)

type Command struct {
	Command   string `json:"command,omitempty"`
	Status    string `json:"status,omitempty"`
	Broadcast string `json:"broadcast,omitempty"`
}

func handleCommand(command *Command, client *Client) {
	if len(command.Command) > 0 {
		fmt.Printf("Received Command: %s\n", command.Command)
		acknowledge := Command{Status: "Acknowledged"}
		client.data <- acknowledge
	}
	if len(command.Status) > 0 {
		fmt.Printf("Status: %s\n", command.Status)
	}
	if len(command.Broadcast) > 0 {
		fmt.Println("Broadcast Message: " + command.Broadcast)
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
