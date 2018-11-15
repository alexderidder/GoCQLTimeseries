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
	NoError        = Error{0, ""}
	FlagNoExist    = Error{0, "Flag doesn't exists"}
	ErrorMarshal   = Error{0, "Problem with parsing/marshall JSON"}
	MissingStoneID = Error{100, "StoneID is missing"}
	MissingType    = Error{100, "Type is missing"}
	MissingData    = Error{100, "Data is missing"}
)

func (e Error) IsNull() bool {
	if e.Code == 0 {
		return true
	}
	return false
}

func (e Error) MarshallErrorAndAddFlag() []byte {
	errBytes := e.marshalError()

	errCode := make([]byte, 4)
	binary.LittleEndian.PutUint32(errCode, 100)
	return append(errCode, errBytes...)
}
func (e Error) marshalError() []byte {
	errBytes, _ := json.Marshal(e)
	//TODO: Marshal error
	return errBytes
}
