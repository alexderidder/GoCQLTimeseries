package model

import "encoding/binary"

type Header struct {
	MessageLength, RequestID, ResponseID, OpCode uint32
}

func (h *Header) makeHeader() []byte {
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

func ByteToArray(request []byte) (Header) {
	result := make([]uint32, 4)
	for i := 0; i < 4; i++ {
		result[i] = ByteToInt(request, i*4)
	}

	return Header{result[0], result[1], result[2], result[3]}
}

func (h *Header) CheckHeader() Error{
	if h.MessageLength == 0 {
		return Error{2, "Header doesn't contain request length"}
	}
	if h.RequestID == 0 {
		return Error{2, "Header doesn't contain request requestID"}
	}
	if h.OpCode == 0 {
	return Error{2, "Header doesn't contain opCode"}
	}
	return Null
}

func ByteToInt(request []byte, beginIndex int) uint32 {
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
