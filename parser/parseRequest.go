package parser

func ParseHeader(request []byte) (uint32, uint32, uint32, uint32) {
	result := make([]uint32, 4)
	for i := 0; i < 4; i++ {
		result[i] = byteToInt(request, i*4)
	}

	return result[0], result[1], result[2], result[3]
}

func ParseOpCode(opCode uint32, message []byte) []byte {
	switch opCode {

	case 100:
		//Insert
	case 200:
		//Select
	default:
		//Return error 'opcode unknown'
	}

	return []byte{}
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
