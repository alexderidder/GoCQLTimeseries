package connector

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
)

//Needs to bigger then 16 (RequestHeader)

type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	socket net.Conn
	data   chan []byte
}

// seems strange to me that the StartServerMode has nothing to do with the server itself. Maybe because server is used for database connection, and this server is a tls server. Improve naming.
func StartServerMode(pem, key, ipAddressAndPort string, timeout, bufferSize uint32) {
	fmt.Println("Starting server...")
	listener := getListenerOverTLS(pem, key, ipAddressAndPort)

	manager := ClientManager{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

	go manager.listenToRegisterAndUnregisterChannelsAndAddOrDelete()

	for {
		connection, err := listener.Accept() // this is blocking, might be nice to explain that in a comment (or is that implied?)
		if err != nil {
			fmt.Println(err) // no logger here?
			continue
		}
		client := &Client{socket: connection, data: make(chan []byte)} // why do we have to call make everywhere?
		manager.register <- client
		go manager.readSocketAndSendToDataChannel(client, bufferSize)
		go client.listenToDataChannelAndProcessMessage(timeout)
	}
}

func getListenerOverTLS(pem, key, ipAddressAndPort string) net.Listener {
	// this seems ...... cryptic. A small comment that would explain this code to readers would be nice.
	cert, err := tls.LoadX509KeyPair(pem, key)

	if err != nil {
		// why do we use a log module here while we're using fmt Println at other places?
		log.Fatal(err)
	}

	tlsConfig := tls.Config{Certificates: []tls.Certificate{cert}, ClientAuth: tls.RequireAnyClientCert}
	tlsConfig.Rand = rand.Reader
	listener, err := tls.Listen("tcp", ipAddressAndPort, &tlsConfig)

	if err != nil {
		fmt.Println(err)
		os.Exit(0) // why not panic?
	}
	return listener
}

func (manager *ClientManager) readSocketAndSendToDataChannel(client *Client, bufferSize uint32) {
	for {
		message := make([]byte, bufferSize)
		length, err := client.socket.Read(message)

		// what sort of errors can happen here that we're forcefully rejecting the client?
		if err != nil {
			manager.unregister <- client
			client.socket.Close()
			break
		}
		if length > 0 {
			client.data <- message[:length]
		}
	}
}
