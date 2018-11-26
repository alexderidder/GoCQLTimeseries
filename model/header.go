package model

import (
	"../util"
	"encoding/binary"
)
const HeaderLength = 16
type Header struct {
	MessageLength, RequestID, ResponseID, OpCode uint32
}

func (h *Header) MakeHeader() []byte {
	var requestHeader []byte
	//Request headers
	variable := make([]byte, 4)

	binary.LittleEndian.PutUint32(variable, h.MessageLength)
	requestHeader = append(requestHeader, variable...)

	binary.LittleEndian.PutUint32(variable, h.RequestID)
	requestHeader = append(requestHeader, variable...)

	binary.LittleEndian.PutUint32(variable, h.ResponseID)
	requestHeader = append(requestHeader, variable...)

	binary.LittleEndian.PutUint32(variable, h.OpCode)
	requestHeader = append(requestHeader, variable...)
	return requestHeader
}

func BytesToHeader(request []byte) (*Header, Error) {
	// no size check, what if this array is not 16 units long? I see the size check in the ByteToUint32 but not here. You will fill your header with false data ("0").
	header := Header{}
	header.MessageLength = util.ByteToUint32(request, 0)
	if header.MessageLength == 0 {

		return nil, HeaderNoLength
	}

	header.RequestID =  util.ByteToUint32(request, 4)
	if header.RequestID == 0 {
		return nil, HeaderNoRequestID
	}
	header.ResponseID =  util.ByteToUint32(request, 8)
	header.OpCode =  util.ByteToUint32(request, 12)
	if header.OpCode == 0 {
		return nil, HeaderNoOpCode
	}

	return &header, NoError
}



