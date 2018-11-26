package util

import "encoding/binary"

func Uint32ToByteArray(value uint32) []byte {
	resultCode := make([]byte, 4) // magic number,
	binary.LittleEndian.PutUint32(resultCode, value)
	return resultCode
}

//Implement error size
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