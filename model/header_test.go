package model

import (
	"GoCQLTimeSeries/util"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)



func TestMakeHeader(t *testing.T) {
	header := Header{1,1,1,1}
	assert.Equal(t, []byte{1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0},  header.MakeHeader(), "Check if header object is correctly transfered to bytes")
	header = Header{128,128,128,128}
	assert.Equal(t, []byte{128, 0, 0, 0, 128, 0, 0, 0, 128, 0, 0, 0, 128, 0, 0, 0},  header.MakeHeader(), "Check if header object is correctly transfered to bytes")

}

func TestByteToArray(t *testing.T) {
	bytes := []byte{1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0}
	header, _ := BytesToHeader(bytes)
	assert.Equal(t, Header{1,1,1,1} , *header, "Check if bytes is correctly transfered is header")
	bytes = []byte{128, 0, 0, 0, 128, 0, 0, 0, 128, 0, 0, 0, 128, 0, 0, 0}
	header, _ = BytesToHeader(bytes)
	assert.Equal(t, Header{128,128,128,128},   *header, "Check if bytes is correctly transfered is header")

}

func TestByteToUint32(t *testing.T) {
	byte1 := []byte{1, 0, 0, 0}
	assert.Equal(t, uint32(1), util.ByteToUint32(byte1, 0), "Test if []byte(1,0,0,0) from index 0 is one")
	byte14336 := []byte{0, 56, 0, 0}
	fmt.Println(util.ByteToUint32(byte14336, 0))

	assert.Equal(t, uint32(14336), util.ByteToUint32(byte14336, 0), "Test if byte [](0,56,0,0) from index 0 is 14336")
	byte128128 := []byte{0, 128, 0, 0, 0}
	assert.Equal(t, uint32(128), util.ByteToUint32(byte128128, 1), "Test if byte [](0,128,0,0,0,0,0) from index 1 is 128")

	byte128128128 := []byte{128, 128, 128}
	assert.Equal(t, uint32(0), util.ByteToUint32(byte128128128, 1), "Test if byte [](128,128,128) from index 0 returns 0")

}