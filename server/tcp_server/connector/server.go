package connector

import (
	"CrownstoneServer/parser"
	"CrownstoneServer/server/config"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"time"
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

func StartServerMode() {

	fmt.Println("Starting server...")
	cert, err := tls.LoadX509KeyPair(config.Config.Server.Certs.Directory+config.Config.Server.Certs.Pem, config.Config.Server.Certs.Directory+config.Config.Server.Certs.Key)

	if err != nil {
		log.Fatal(err)
	}

	tlsConfig := tls.Config{Certificates: []tls.Certificate{cert}, ClientAuth: tls.RequireAnyClientCert}
	tlsConfig.Rand = rand.Reader
	listener, err := tls.Listen("tcp", config.Config.Server.IPAddress+config.Config.Server.Port, &tlsConfig)

	if err != nil {
		fmt.Println(err)
	}
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
		go manager.receive(client)
		go manager.processMessage(client)
	}
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

func (manager *ClientManager) receive(client *Client) {
	for {
		message := make([]byte, config.Config.Server.Messages.BufferSize)
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

func (manager *ClientManager) processMessage(client *Client) {
	defer client.socket.Close()
	for {
		if !client.processReceivedMessage() {
			return
		}
	}
}

func (client *Client) processReceivedMessage() bool {
	var result []byte
	select {
	case message, ok := <-client.data:
		if !ok {
			return false
		}

		requestLength, requestID, _, opCode := parser.ParseHeader(message)
		result = message
		if requestLength == 0 {
			client.sendResponseError(2, "Header doesn't contain request length", 0)
		} else if requestID == 0 {
			client.sendResponseError(2, "Header doesn't contain request ID", 0)
		} else if opCode == 0 {
			client.sendResponseError(2, "Header doesn't contain opCode", 0)
		} else {
			for requestLength > uint32(len(message)) {
				select {
				case custom, ok := <-client.data:
					if !ok {
						return false
					} else {
						result = append(result, custom...)
					}
				case <-time.After(time.Duration(config.Config.Server.Messages.Timeout) * time.Second):
					client.sendResponseError(20, "Server didn't receive full message", requestID)
					return true
				}

			}
			if len(result) > int(requestLength) {
				result = result[:requestLength]
			}
			result = result[16:]
			responseWithoutHeader := parser.ParseOpCode(opCode, result)
			client.sendResponseWithoutHeader(requestID, responseWithoutHeader)

		}
	}
	return true
}

func (client *Client) sendResponseWithoutHeader(requestID uint32, responseWithoutHeader []byte) {
	//fmt.Println(uint32(len(responseWithoutHeader)) + 16)
	request := append(parser.MakeHeader(uint32(len(responseWithoutHeader)+16), 0, requestID, 1), responseWithoutHeader...)
	_, err := client.socket.Write(request)
	if err != nil {
		//TODO: write error
		fmt.Println(-1)
	}
}

func (client *Client) sendResponseError(code uint32, message string, requestID uint32) {
	responseWithoutHeader := parser.ParseError(2, "Header doesn't contain opCode")
	client.sendResponseWithoutHeader(requestID, responseWithoutHeader)
}
