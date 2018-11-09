package connector

import (
	"CrownstoneServer/model"
	"CrownstoneServer/parser"
	"encoding/binary"
	"fmt"
	"time"
)

func (client *Client) listenToDataChannelAndProcessMessage(timeout uint32)  {
	defer client.socket.Close()
	var result []byte
	for {
		select {
		case message, ok := <-client.data:
			if !ok {
				return
			}

			request := model.ByteToArray(message)
			err := request.CheckHeader()
			if !err.IsNull() {
				fmt.Println(err)
				client.writeJsonResponse(request.RequestID, err.MarshallErrorAndAddFlag())
			}
			result = message


			for request.MessageLength > uint32(len(result)) {
				select {
				case custom, ok := <-client.data:
					if !ok {
						return
					} else {
						result = append(result, custom...)
					}
				case <-time.After(time.Duration(timeout) * time.Second):
					errBytes := model.Error{20, "Server didn't receive full message"}.MarshallErrorAndAddFlag()
					client.writeJsonResponse(request.RequestID, errBytes)
					continue
				}

			}
			if len(result) > int(request.MessageLength) {
				result = result[:request.MessageLength]
			}
			result = result[16:]
			responseWithoutHeader := parser.ParseOpCode(request.OpCode, result)
			client.writeJsonResponse(request.RequestID, responseWithoutHeader)
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