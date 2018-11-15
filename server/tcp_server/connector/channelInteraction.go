package connector

import (
	"CrownstoneServer/model"
	"CrownstoneServer/parser"
	"fmt"
	"time"
)

func (manager *ClientManager) listenToRegisterAndUnregisterChannelsAndAddOrDelete() {
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

func (client *Client) listenToDataChannelAndProcessMessage(timeout uint32)  {
	defer client.socket.Close()
	var result []byte
	for {
		select {
		case message, ok := <-client.data:
			if !ok {
				return
			}
			requestHeader := model.BytesToHeader(message)
			error := requestHeader.CheckHeader()
			if !error.IsNull() {
				client.writeJsonResponse(requestHeader.RequestID, error.MarshallErrorAndAddFlag())
				continue
			}
			 result = message[16:]
			//When expected length is smaller then the sum of the received messages -> read again
			for requestHeader.MessageLength - 16 > uint32(len(result)) {
				select {
				case custom, ok := <-client.data:
					if !ok {
						return
					} else {
						result = append(result, custom...)
					}
				case <-time.After(time.Duration(timeout) * time.Second):
					errBytes := model.Error{20, "Server didn't receive full message"}.MarshallErrorAndAddFlag()
					client.writeJsonResponse(requestHeader.RequestID, errBytes)
					continue
				}
			}

			responseWithoutHeader := parser.ProcessOpCodeAndReceivedMessage(requestHeader.OpCode, result)
			client.writeJsonResponse(requestHeader.RequestID, responseWithoutHeader)
		}
	}

}



func (client *Client) writeJsonResponse(requestID uint32, responseWithoutHeader []byte) {
	header := model.Header{uint32(len(responseWithoutHeader)+16), 0, requestID, 1}

	request := append(header.MakeHeader(), responseWithoutHeader...)
	_, err := client.socket.Write(request)
	if err != nil {
		//TODO: write error
		fmt.Println(-1)
	}
}