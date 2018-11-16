package model

import "encoding/binary"

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

func BytesToHeader(request []byte) Header {
	// no size check, what if this array is not 16 units long? I see the size check in the ByteToUint32 but not here. You will fill your header with false data ("0").
	// the meaning of number 0 as error is not normal behaviour and is not commented on anywhere.
	result := make([]uint32, 4)
	for i := 0; i < 4; i++ {
		result[i] = ByteToUint32(request, i*4)
	}

	// this might be filled with 0's, which is a valid parsing of bytes AND it's your definition of error. This is unclear.
	return Header{result[0], result[1], result[2], result[3]}
}

func (h *Header) CheckHeader() Error {
	// this can be combined with the constructor to return the consistant format (instance, err := ...
	if h.MessageLength == 0 {
		return Error{2, "Header doesn't contain request length"}
	}
	if h.RequestID == 0 {
		return Error{2, "Header doesn't contain request requestID"}
	}
	if h.OpCode == 0 {
		return Error{2, "Header doesn't contain opCode"}
	}
	return NoError
}

func ByteToUint32(request []byte, beginIndex int) uint32 {
	if len(request) >= beginIndex+4 {
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
	return 0

}
