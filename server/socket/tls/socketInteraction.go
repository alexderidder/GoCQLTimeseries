package tls

import (
	"GoCQLTimeSeries/model"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"time"
)

type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	socket net.Conn
}

func StartTLSServer(pem, key, ipAddressAndPort string, timeout, bufferSize uint32) {
	fmt.Println("Starting server...")
	listener := getListenerOverTLS(pem, key, ipAddressAndPort)

	manager := ClientManager{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

	go manager.listenToRegisterAndUnregisterChannelsAndAddOrDelete()

	for {
		// Accept waits for and returns the next connection to the listener.
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		client := &Client{socket: connection}
		manager.register <- client
		go manager.readSocketAndSendToDataChannel(client)
	}
}

func getListenerOverTLS(cert, key, ipAddressAndPort string) net.Listener {
	//Read the cert and key files and save them as a Certificate
	certificate, err := tls.LoadX509KeyPair(cert, key)

	if err != nil {
		//Close program because it needs cert for Encryption
		panic(err)
	}
	//tlsConfig := tls.Config{Certificates: []tls.Certificate{certificate}, ClientAuth: tls.VerifyClientCertIfGiven}
	tlsConfig := tls.Config{Certificates: []tls.Certificate{certificate}, ClientAuth: tls.RequireAnyClientCert}

	listener, err := tls.Listen("tcp", ipAddressAndPort, &tlsConfig)
	if err != nil {
		//Close program because it needs the right config values
		panic(err)
	}
	return listener
}

func (manager *ClientManager) listenToRegisterAndUnregisterChannelsAndAddOrDelete() {
	for {
		select {
		case connection := <-manager.register:
			manager.clients[connection] = true
			fmt.Println("Added new connection!")
		case connection := <-manager.unregister:
			if _, ok := manager.clients[connection]; ok {
				delete(manager.clients, connection)
				fmt.Println("A connection has terminated!")
			}
		}
	}
}

func (manager *ClientManager) readSocketAndSendToDataChannel(client *Client) {

	headerBytes := make([]byte, model.HeaderLength)
	var message []byte

	//Keeps reading messages, first header. Then the remaining bytes.
	for {
		client.socket.SetReadDeadline(time.Time{})
		err := client.read(&headerBytes)
		if err != nil {
			//Can only receive connection close message
			manager.unregister <- client
			client.socket.Close()
			break
		}

		//Parse header
		requestHeader, error := model.BytesToHeader(headerBytes)
		if !error.IsNull() {
			fmt.Println(error.Message)
			client.writeJsonResponse(requestHeader.RequestID, error.MarshallErrorAndAddFlag())
			continue

		}
		//Read payload with length declared in header - headerlength
		message = make([]byte, int(requestHeader.MessageLength)-model.HeaderLength)
		client.socket.SetReadDeadline(time.Now().Add(time.Millisecond * 100))
		err = client.read(&message)
		if err != nil {
			//if err is timeout, payload is smaller then messagelength. So return error
			if err, ok := err.(net.Error); ok && err.Timeout() {
				{
					errBytes := model.ReceivedFullMessage.MarshallErrorAndAddFlag()
					client.writeJsonResponse(requestHeader.RequestID, errBytes)
					continue
				}

			}
			//Else it is connect close error message.
			manager.unregister <- client
			client.socket.Close()
			break;
		}
		//If received header + payload. Process message
		go client.parseExecuteAndResponseToMessage(*requestHeader, message)

	}

}

func (client *Client) read(bytes *[]byte) error {
	//Returns error if receives connection close message or timeout
	_, err := io.ReadFull(client.socket, *bytes)
	return err
}
