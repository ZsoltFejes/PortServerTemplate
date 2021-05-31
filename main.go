package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var WORKDIR string
var debug bool

// Check Error function for universal error handling
func checkErr(message string, err error) {
	if err != nil {
		l(message+"- "+err.Error(), true, true)
	}
}

func l(message string, fatal bool, public bool) {
	if public || debug {
		fmt.Println(message)
	}
	if fatal {
		log.Fatalln(message)
	} else {
		log.Println(message)
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
	flagMode := flag.String("mode", "client", "Start in client or server mode")
	flagTLS := flag.Bool("tls", false, "Set Server to use TLS (Add certifiacet to root directory as server.crt and server.key)")
	flag.BoolVar(&debug, "debug", false, "Set process to debug")
	flag.Parse()

	// Set working directory
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	WORKDIR = filepath.Dir(ex)

	// Create or open log directory
	f, err := os.OpenFile(WORKDIR+`/server.log`, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		l(err.Error(), true, true)
	}
	defer f.Close()
	log.SetOutput(f)

	// Start application in requested mode
	if strings.ToLower(*flagMode) == "server" {
		startServerMode(&server, flagTLS)
	} else if strings.ToLower(*flagMode) == "client" {
		startClientMode(flagTLS)
	} else {
		fmt.Println("Mode is unknown!")
	}
}
