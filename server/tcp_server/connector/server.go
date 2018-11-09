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

type Error struct {
	ErrorMessage string
}

type Client struct {
	socket net.Conn
	data   chan []byte

}

func StartServerMode(pem , key, ipaddressAndPort string, timeout, bufferSize uint32 ){

	fmt.Println("Starting server...")
	listener := getListenerOverTLS(pem, key, ipaddressAndPort)

	manager := ClientManager{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	go manager.start()

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		client := &Client{socket: connection, data: make(chan []byte)}
		manager.register <- client
		go manager.receive(client, bufferSize)
		go client.listenToDataChannelAndProcessMessage(timeout)
	}
}

func getListenerOverTLS(pem, key, ipaddressAndPort string ) net.Listener{
	cert, err := tls.LoadX509KeyPair(pem, key)

	if err != nil {
		log.Fatal(err)
	}

	tlsConfig := tls.Config{Certificates: []tls.Certificate{cert}, ClientAuth: tls.RequireAnyClientCert}
	tlsConfig.Rand = rand.Reader
	listener, err := tls.Listen("tcp", ipaddressAndPort, &tlsConfig)

	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	return listener
}

func (manager *ClientManager) start() {
	for {
		select {
		case connection := <-manager.register:
			manager.clients[connection] = true
			fmt.Println("Added new connection!")
		case connection := <-manager.unregister:
			if _, ok := manager.clients[connection]; ok {
				close(connection.data)
				delete(manager.clients, connection)
				fmt.Println("A connection has terminated!")
			}
		}
	}
}

func (manager *ClientManager) receive(client *Client, bufferSize uint32) {
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


