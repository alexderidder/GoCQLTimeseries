package _select

import (
	"GoCQLTimeSeries/datatypes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelectFlag1JSON(t *testing.T) {
	message := []byte("{\"stoneIDs\":[\"test123\"],\"startTime\":1548258873000,\"endTime\":1548258873000,\"interval\":0}")
	_, err := parseFlag1(message, 0)
	assert.Equal(t, datatypes.NoError, err, "Check if no Error")
	message = []byte("")
	_, err = parseFlag1(message, 0)
	assert.Equal(t, datatypes.UnMarshallError, err, "Check for marshal error unexpected end")

	message = []byte("{\"stoneIDs\":[],\"startTime\":1548258873000,\"endTime\":1548258873000,\"interval\":0}")
	_, err = parseFlag1(message, 0)
	assert.Equal(t, datatypes.MissingStoneID, err, "Check for missing stone id")


	message = []byte("{\"stoneIDs\":[\"test123\"],\"interval\":0}")
	_, err = parseFlag1(message, 0)
	assert.Equal(t, datatypes.MissingStartAndEndTime, err, "Check for missing start and end time")

	message = []byte("{\"stoneIDs\":[\"test123\"],\"endTime\":1548258873000,\"interval\":0}")
	_, err = parseFlag1(message, 0)
	assert.Equal(t, datatypes.MissingStartTime, err, "Check for missing start time")

	message = []byte("{\"stoneIDs\":[\"test123\"],\"startTime\":1548258873000,\"interval\":0}")
	_, err = parseFlag1(message, 0)
	assert.Equal(t, datatypes.MissingEndTime, err, "Check for missing end time")
}
