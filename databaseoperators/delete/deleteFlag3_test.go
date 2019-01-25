package delete

import (
	"GoCQLTimeSeries/datatypes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleteFlag3JSON(t *testing.T) {
	message := []byte("{\"stoneID\":[\"test1234\"]}")
	_, err := parseFlag3(message, 0)
	assert.Equal(t, datatypes.NoError, err, "Check if no Error")
	message = []byte("")
	_, err = parseFlag3(message, 0)
	assert.Equal(t, datatypes.Error{301, "unexpected end of JSON input"}, err, "Check for marshal error unexpected end")

	message = []byte("{\"stoneID\":[]}")
	_, err = parseFlag3(message, 0)
	assert.Equal(t, datatypes.MissingStoneID, err, "Check for missing stone id")


}