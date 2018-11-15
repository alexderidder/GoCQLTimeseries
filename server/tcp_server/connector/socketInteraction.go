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


func StartServerMode(pem , key, ipAddressAndPort string, timeout, bufferSize uint32 ){
	fmt.Println("Starting server...")
	listener := getListenerOverTLS(pem, key, ipAddressAndPort)

	manager := ClientManager{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

	go manager.listenToRegisterAndUnregisterChannelsAndAddOrDelete()

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		client := &Client{socket: connection, data: make(chan []byte)}
		manager.register <- client
		go manager.readSocketAndSendToDataChannel(client, bufferSize)
		go client.listenToDataChannelAndProcessMessage(timeout)
	}
}

func getListenerOverTLS(pem, key, ipAddressAndPort string ) net.Listener{
	cert, err := tls.LoadX509KeyPair(pem, key)

	if err != nil {
		log.Fatal(err)
	}

	tlsConfig := tls.Config{Certificates: []tls.Certificate{cert}, ClientAuth: tls.RequireAnyClientCert}
	tlsConfig.Rand = rand.Reader
	listener, err := tls.Listen("tcp", ipAddressAndPort, &tlsConfig)

	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	return listener
}



func (manager *ClientManager) readSocketAndSendToDataChannel(client *Client, bufferSize uint32) {
	for {
		message := make([]byte, bufferSize)
		length, err := client.socket.Read(message)
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


