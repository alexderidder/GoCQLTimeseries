package insert

import (
	"GoCQLTimeSeries/datatypes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInsertFlag2JSON(t *testing.T) {
	message := []byte("{\"stoneID\":\"test1234\", \"data\":[{\"time\":1548258873000,\"value\" : {\"kWh\":1.5210}}]}")
	_, err := parseFlag2(message, 0)
	assert.Equal(t, datatypes.NoError, err, "Check if no Error")
	message = []byte("")
	_, err = parseFlag2(message, 0)
	assert.Equal(t, datatypes.UnMarshallError, err, "Check for marshal error unexpected end")

	message = []byte("{\"data\":[{\"time\":1548258873000,\"value\" : {\"kWh\":1.5210}}]}")
	_, err = parseFlag2(message, 0)
	assert.Equal(t, datatypes.MissingStoneID, err, "Check for missing stone id")
	message = []byte("{\"stoneID\":\"test1234\"}")
	_, err = parseFlag2(message, 0)
	assert.Equal(t, datatypes.MissingData, err, "Check if data is  missing")
}
