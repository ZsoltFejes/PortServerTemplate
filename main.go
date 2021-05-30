package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Check Error function for universal error handling
func checkErr(message string, err interface{}) {
	if err != nil {
		fmt.Printf("%s- %s\n", message, err)
	}
}

// Random string generator to identify Jobs
func getID() string {
	var characters = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, 8)
	for i := range b {
		b[i] = characters[rand.Intn(len(characters))]
	}
	return string(b)
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Seed random with current time
	flagMode := flag.String("mode", "server", "Start in client or server mode")
	flagTLS := flag.Bool("tls", false, "Set Server to use TLS (Add certifiacet to root directory as server.crt and server.key)")
	flag.Parse()
	if strings.ToLower(*flagMode) == "server" {
		startServerMode(&server, flagTLS)
	} else if strings.ToLower(*flagMode) == "client" {
		startClientMode(flagTLS)
	} else {
		fmt.Println("Mode is unknown!")
	}
}
