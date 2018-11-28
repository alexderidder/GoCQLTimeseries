package tls

import (
	"GoCQLTimeSeries/model"
	"GoCQLTimeSeries/processMessage"
	"fmt"
)

func (client *Client) parseExecuteAndResponseToMessage(requestHeader model.Header, result []byte) {

	executeObject, err := processMessage.ParseOpCode(requestHeader.OpCode, &result)
	if !err.IsNull() {
		client.writeJsonResponse(requestHeader.RequestID, err.MarshallErrorAndAddFlag())
		return
	}
	responsePayload, err := executeObject.Execute()
	if !err.IsNull() {
		client.writeJsonResponse(requestHeader.RequestID, err.MarshallErrorAndAddFlag())
		return
	}


	client.writeJsonResponse(requestHeader.RequestID, responsePayload)

}

func (client *Client) writeJsonResponse(requestID uint32, responseWithoutHeader []byte) {
	header := model.Header{uint32(len(responseWithoutHeader)) + model.HeaderLength, 0, requestID, 1}

	request := append(header.MakeHeader(), responseWithoutHeader...)
	_, err := client.socket.Write(request)
	fmt.Println(request)
	if err != nil {
		//When connection is closed error is returned. This goroutine will end after sending the message. So no further actions needed
		fmt.Println(err)
	}
}
