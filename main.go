package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Server PortIP `json:"server,omitempty"`
	Api    PortIP `json:"api,omitempty"`
	Client PortIP `json:"client,omitempty"`
	Tls    bool   `json:"tls"`
}

type PortIP struct {
	Address string `json:"address,omitempty"`
	Port    string `json:"port,omitempty"`
}

var (
	WORKDIR    string
	debug      = flag.Bool("debug", false, "Set process to debug")
	configFile = flag.String("config", "config.json", "Specify the location of the config file.")
	appConfig  Config
)

// Check Error function for universal error handling
func checkErr(message string, err error, showError bool) {
	if err != nil {
		if showError {
			l(message+"- "+err.Error(), true, showError)
		} else {
			l(message, false, !showError)            // Show error without the error
			l("ERROR "+err.Error(), true, showError) // Log error with the error
		}
	}
}

func l(message string, fatal bool, public bool) {
	if public || *debug {
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
	flag.Parse()

	// Set working directory
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	WORKDIR = filepath.Dir(ex)

	// Load config file
	f, err := ioutil.ReadFile(WORKDIR + *configFile)
	checkErr("Unable to read configuration file, check the path to the file", err, false)
	err = json.Unmarshal(f, &appConfig)
	checkErr("Unable to load config file. Please check the documentation", err, false)

	// Start application in requested mode
	if len(appConfig.Server.Port) > 0 {
		startServerMode()
	} else if len(appConfig.Client.Port) > 0 {
		startClientMode()
	} else {
		fmt.Println("Mode is unknown!")
	}
}
