package main

import (
	"flag"
	"fmt"
	"strings"
)

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
