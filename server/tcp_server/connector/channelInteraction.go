package connector

import (
	"../../../model"
	"../../../parser"
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

func (client *Client) listenToDataChannelAndProcessMessage(timeout uint32) {
	defer client.socket.Close()
	var result []byte
	for {
		select {
		case message, ok := <-client.data:
			if !ok { // why OK and not err like everywhere else? Where do we return to? Is the goroutine then terminated?
				return
			}
			requestHeader := model.BytesToHeader(message)
			// to be consistent with the other methods, when creating a header, also return an error (requestHeader, err := ...)
			error := requestHeader.CheckHeader() // get rid of this.
			if !error.IsNull() {
				client.writeJsonResponse(requestHeader.RequestID, error.MarshallErrorAndAddFlag())
				continue
			}

			// what is result? Is it message body?
			result = message[16:] // where does this magic 16 come from? define HEADER_LENGTH somewhere.
			//When expected length is smaller then the sum of the received messages -> read again
			for requestHeader.MessageLength-16 > uint32(len(result)) { // where does this magic 16 come from? define HEADER_LENGTH somewhere.
				// I don't follow this. What's happening here? You already have the message, why are we waiting on more data from the client?
				select {
				case custom, ok := <-client.data: // what is a custom?
					if !ok {
						return
					} else {
						result = append(result, custom...) // surely we can be cleaner than copying one byte at a time...
					}
				case <-time.After(time.Duration(timeout) * time.Second):
					errBytes := model.Error{20, "Server didn't receive full message"}.MarshallErrorAndAddFlag() // error code 20 is not defined in protocol doc.
					client.writeJsonResponse(requestHeader.RequestID, errBytes)
					continue
				}
			}

			responseWithoutHeader := parser.ProcessOpCodeAndReceivedMessage(requestHeader.OpCode, result)

			// comment what is happening here. This is a protocol agreement (empty response) but does it allow for async processing of the message?
			client.writeJsonResponse(requestHeader.RequestID, responseWithoutHeader)
		}
	}

}

func (client *Client) writeJsonResponse(requestID uint32, responseWithoutHeader []byte) {
	header := model.Header{uint32(len(responseWithoutHeader) + 16), 0, requestID, 1}

	request := append(header.MakeHeader(), responseWithoutHeader...)
	_, err := client.socket.Write(request)
	if err != nil {
		//TODO: write error
		fmt.Println(-1)
	}
}
