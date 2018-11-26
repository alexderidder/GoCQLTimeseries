package model

import (
	"encoding/binary"
	"encoding/json"
)

type Error struct {
	Code    uint32 `json:"code"`
	Message string `json:"message"`
}

var (
	NoError             = Error{0, ""}
	FlagNoExist         = Error{150, "Flag doesn't exists"}
	ErrorMarshal        = Error{150, "Problem with parsing/marshall JSON"}
	MissingStoneID      = Error{100, "StoneID is missing"}
	MissingType         = Error{100, "Type is missing"}
	MissingData         = Error{100, "Data is missing"}
	ReceivedFullMessage = Error{20, "Server didn't receive full message"}
	ServerNoCassandra = Error{21, "Can't connect to cassandra"}
	HeaderNoLength      = Error{2, "Header doesn't contain request length"}
	HeaderNoRequestID   = Error{2, "Header doesn't contain request requestID"}
	HeaderNoOpCode      = Error{2, "Header doesn't contain opCode"}
	MarshallError       = Error{300, ""}
	UnMarshallError     = Error{301, ""}
)

func (e *Error) IsNull() bool {
	if e.Code == 0 {
		return true
	}
	return false
}

func (e *Error) MarshallErrorAndAddFlag() []byte {
	errBytes := e.marshalError()

	errCode := make([]byte, 4)
	binary.LittleEndian.PutUint32(errCode, 100)
	return append(errCode, errBytes...)
}
func (e *Error) marshalError() []byte {
	errBytes, _ := json.Marshal(e)
	//TODO: Marshal error
	return errBytes
}
