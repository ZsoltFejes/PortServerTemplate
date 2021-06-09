package main

import (
	"crypto/rand"
	"crypto/tls"
	"log"
	"net"
	"os"
	"time"
)

type Server struct {
	clients    map[*Client]bool
	broadcast  chan Job
	register   chan *Client
	unregister chan *Client
}

var server = Server{
	clients:    make(map[*Client]bool),
	broadcast:  make(chan Job),
	register:   make(chan *Client),
	unregister: make(chan *Client),
}

// Start Client manager
func (server *Server) start() {
	l("Server started "+appConfig.Server.Address+":"+appConfig.Server.Port, false, true)
	for {
		select {
		case connection := <-server.register:
			server.clients[connection] = true
			l("Added new connection! "+connection.socket.RemoteAddr().String(), false, false)
		case connection := <-server.unregister:
			if _, ok := server.clients[connection]; ok {
				close(connection.data)
				delete(server.clients, connection)
				l("A connection has been terminated "+connection.socket.RemoteAddr().String(), false, false)
			}
		case message := <-server.broadcast:
			for connection := range server.clients {
				select {
				case connection.data <- message:
				default:
					close(connection.data)
					delete(server.clients, connection)
				}
			}
		}
	}
}

/* Starts server specific goroutines.
The function creates a TLS or unenctypted socket based on the configuraiton. Then starts listening on the specified port for incoming TCP requests.
If a TCP connection is esablished, the connection will be registered as aclient then the server starts a go routine for receiving data from and for sending data to the client.
*/
func startServerMode() {
	// Create or open log directory
	f, err := os.OpenFile(WORKDIR+`/server.log`, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		l(err.Error(), true, true)
	}
	defer f.Close()
	log.SetOutput(f)
	l("Starting server...", false, true)
	var listener net.Listener
	if appConfig.Tls {
		cert, err := tls.LoadX509KeyPair(WORKDIR+"/cert.pem", WORKDIR+"/key.pem")
		checkErr("Importing TLS certificates error", err)
		config := tls.Config{Certificates: []tls.Certificate{cert}}
		now := time.Now()
		config.Time = func() time.Time { return now }
		config.Rand = rand.Reader
		listener, err = tls.Listen("tcp", appConfig.Server.Address+":"+appConfig.Server.Port, &config)
		checkErr("Creating TLS listener error", err)
	} else {
		var err error
		listener, err = net.Listen("tcp", appConfig.Server.Address+":"+appConfig.Server.Port)
		checkErr("Creating NET listener error", err)
	}
	go server.start()
	if len(appConfig.Api.Port) > 0 {
		go startHttpServer()
	}
	for {
		connection, err := listener.Accept()
		checkErr("Accepting connection error", err)
		client := &Client{socket: connection, data: make(chan Job)}
		server.register <- client
		go client.receive()
		go client.send()
	}
}
