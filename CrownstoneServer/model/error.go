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
	Null  = Error{0, ""}
	MissingStoneID = Error{100, "StoneID is missing"}
	MissingType = Error{100, "Type is missing"}
)

func (e Error) IsNull() bool{
	if e.Code == 0 {
		return true
	}
	return false
}
func (e Error) checkError() {

	if e.Code == 0 && len(e.Message) == 0 {
		e = Error{500, "Server created error without code and message"}
	} else if e.Code == 0 {
		e = Error{500, "Server created error without code and message"}
	} else if len(e.Message) == 0 {
		e = Error{500, "Server created error without code and message"}
	}
}
func (e Error) MarshallErrorAndAddFlag() []byte{
	errBytes := e.marshalError()

	errCode := make([]byte, 4)
	binary.LittleEndian.PutUint32(errCode, 100)
	return append(errCode, errBytes...)
}
func (e Error) marshalError() []byte{
	errBytes, _ := json.Marshal(e)
	//TODO: Marshal error
	return errBytes
}



