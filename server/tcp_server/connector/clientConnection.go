package connector

import (
	"CrownstoneServer/parser"
	"CrownstoneServer/server/config"
	"encoding/binary"
	"fmt"
	"time"
)

func (client *Client) listenToDataChannelAndProcessMessage()  {
	defer client.socket.Close()
	var result []byte
	for {
		select {
		case message, ok := <-client.data:
			if !ok {
				return
			}

			requestLength, requestID, _, opCode := parser.ParseHeader(message)
			result = message
			if requestLength == 0 {
				errBytes := parser.ParseError(2, "Header doesn't contain request length")
				client.writeJsonResponse(0, errBytes)
				continue
			}
			if requestID == 0 {
				errBytes := parser.ParseError(2, "Header doesn't contain request requestID")
				client.writeJsonResponse(0, errBytes)
				continue
			}
			if opCode == 0 {
				errBytes := parser.ParseError(2, "Header doesn't contain opCode")
				client.writeJsonResponse(requestID, errBytes)
				continue
			}

			for requestLength > uint32(len(result)) {
				select {
				case custom, ok := <-client.data:
					if !ok {
						return
					} else {
						result = append(result, custom...)
					}
				case <-time.After(time.Duration(config.Config.Server.Messages.Timeout) * time.Second):
					errBytes := parser.ParseError(20, "Server didn't receive full message")
					client.writeJsonResponse(requestID, errBytes)
					continue
				}

			}
			if len(result) > int(requestLength) {
				result = result[:requestLength]
			}
			result = result[16:]
			responseWithoutHeader := parser.ParseOpCode(opCode, result)
			client.writeJsonResponse(requestID, responseWithoutHeader)
		}
	}

}


func (client *Client) writeJsonResponse(requestID uint32, responseWithoutHeader []byte) {
	//fmt.Println(uint32(len(responseWithoutHeader)) + 16)
	request := append(client.makeHeader(uint32(len(responseWithoutHeader)+16), 0, requestID, 1), responseWithoutHeader...)
	_, err := client.socket.Write(request)
	if err != nil {
		//TODO: write error
		fmt.Println(-1)
	}
}

func (client *Client) makeHeader(messageLength, requestID, responseID, opCode uint32) []byte {
	var requestHeader []byte
	//Request headers
	variable := make([]byte, 4)

	binary.LittleEndian.PutUint32(variable, messageLength)
	requestHeader = append(requestHeader, variable...)

	binary.LittleEndian.PutUint32(variable, requestID)
	requestHeader = append(requestHeader, variable...)

	binary.LittleEndian.PutUint32(variable, responseID)
	requestHeader = append(requestHeader, variable...)

	binary.LittleEndian.PutUint32(variable, opCode)
	requestHeader = append(requestHeader, variable...)
	return requestHeader
}