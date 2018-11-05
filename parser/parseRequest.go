package parser

import (
	"CrownstoneServer/server/database/connector"
	"encoding/binary"
	"encoding/json"
)

type Error struct {
	Code    uint32 `json:"code"`
	Message string `json:"message"`
}

func ParseHeader(request []byte) (uint32, uint32, uint32, uint32) {
	result := make([]uint32, 4)
	for i := 0; i < 4; i++ {
		result[i] = byteToInt(request, i*4)
	}

	return result[0], result[1], result[2], result[3]
}

func ParseOpCode(opCode uint32, message []byte) []byte {
	if cassandra.Session == nil{
		return ParseError(300, "There is no connection between GoCQL-driver and Cassandra")
	}
	switch opCode {
	//TODO: insert
	//case 100:
	//TODO: select
	case 200:
		return parseSelect(message)
	default:
		return ParseError(10, "Server doesn't recognise opcode")
	}
}

func byteToInt(request []byte, beginIndex int) uint32 {
	var result uint32
	result |= uint32(request[beginIndex])
	beginIndex++
	result |= uint32(request[beginIndex]) << 8
	beginIndex++
	result |= uint32(request[beginIndex]) << 16
	beginIndex++
	result |= uint32(request[beginIndex]) << 24
	return result
}

func ParseError(code uint32, message string) []byte {
	var error Error
	if code == 0 && len(message) == 0 {
		error = Error{500, "Server created error without code and message"}
	} else if code == 0 {
		error = Error{500, "Server created error without code"}
	} else if len(message) == 0 {
		error = Error{500, "Server created error without message"}
	} else {
		error = Error{code, message}
	}

	errBytes, _ := json.Marshal(error)
	//TODO: Marshal error
	errCode := make([]byte, 4)
	binary.LittleEndian.PutUint32(errCode, 100)
	return append(errCode, errBytes...)
}

func MakeHeader(messageLength, requestID, responseID, opCode uint32) []byte {
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
