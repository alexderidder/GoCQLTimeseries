package tls

import (
	"GoCQLTimeSeries/model"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"net"
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
	tlsConfig, err := createServerConfig("server/socket/certs/ca-crt.pem", cert, key)
	if err != nil {
		panic( err.Error())
	}

	listener, err := tls.Listen("tcp", ipAddressAndPort, tlsConfig)
	if err != nil {
		//Close program because it needs the right config values
		panic(err)
	}
	return listener
}

func createServerConfig(ca, crt, key string) (*tls.Config, error) {
	caCertPEM, err := ioutil.ReadFile(ca)
	if err != nil {
		return nil, err
	}

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(caCertPEM)
	if !ok {
		panic("failed to parse root certificate")
	}

	cert, err := tls.LoadX509KeyPair(crt, key)
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    roots,
	}, nil
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
		//client.socket.SetReadDeadline(time.Time{})
		_, err := client.read(headerBytes)
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
			manager.unregister <- client
			client.socket.Close()
			return
		}
		//Read payload with length declared in header - headerlength
		message = make([]byte, int(requestHeader.MessageLength)-model.HeaderLength)
		//client.socket.SetReadDeadline(time.Now().Add(time.Millisecond * 300))
		_, err = client.read(message)
		if err != nil {
			//if err is timeout, payload is smaller then messagelength. So return error
			//if err, ok := err.(net.Error); ok && err.Timeout() {
			//	{
			//		errBytes := model.ReceivedFullMessage.MarshallErrorAndAddFlag()
			//		client.writeJsonResponse(requestHeader.RequestID, errBytes)
			//		manager.unregister <- client
			//		client.socket.Close()
			//		return
			//	}
			//
			//}
			//Else it is connect close error message.
			manager.unregister <- client
			client.socket.Close()
			break
		}
		//If received header + payload. Process message
		go client.parseExecuteAndResponseToMessage(*requestHeader, message)

	}

}

func (client *Client) read(bytes []byte) (int, error) {
	//Returns error if receives connection close message or timeout
	return io.ReadFull(client.socket, bytes)
}
