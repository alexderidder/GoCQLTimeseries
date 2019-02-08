package insert

import (
	"GoCQLTimeSeries/datatypes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInsertFlag1JSON(t *testing.T) {
	message := []byte("{\"stoneID\":\"test1234\", \"data\":[{\"time\":1548258873000,\"value\" : {\"watt\":3.0233014,\"pf\":4.702545}}]}")
	_, err := parseFlag1(message, 0)
	assert.Equal(t, datatypes.NoError, err, "Check if no Error")
		message = []byte("")
		_, err = parseFlag1(message, 0)
		assert.Equal(t, datatypes.UnMarshallError, err, "Check for marshal error unexpected end")

		message = []byte("{\"data\":[{\"time\":1548258873000,\"value\" : {\"watt\":3.0233014,\"pf\":4.702545}}]}")
		_, err = parseFlag1(message, 0)
		assert.Equal(t, datatypes.MissingStoneID, err, "Check for missing stone id")
		message = []byte("{\"stoneID\":\"test1234\"}")
		_, err = parseFlag1(message, 0)
		assert.Equal(t, datatypes.MissingData, err, "Check if data is  missing")
}
