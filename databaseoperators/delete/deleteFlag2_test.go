package delete

import (
	"GoCQLTimeSeries/datatypes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleteFlag2JSON(t *testing.T) {
	message := []byte("{\"stoneID\":\"test1234\",\"startTime\":1548258873000,\"endTime\":1548258873000}")
	_, err := parseFlag2(message, 0)
	assert.Equal(t, datatypes.NoError, err, "Check if no Error")
	message = []byte("")
	_, err = parseFlag2(message, 0)
	assert.Equal(t, datatypes.UnMarshallError, err, "Check for marshal error unexpected end")

	message = []byte("{\"startTime\":1548258873000,\"endTime\":1548258873000}")
	_, err = parseFlag2(message, 0)
	assert.Equal(t, datatypes.MissingStoneID, err, "Check for missing stone id")

	message = []byte("{\"stoneID\":\"test1234\"}")
	_, err = parseFlag2(message, 0)
	assert.Equal(t, datatypes.MissingStartAndEndTime, err, "Check for missing start and end time")

	message = []byte("{\"stoneID\":\"test1234\",\"endTime\":1548258873000}")
	_, err = parseFlag2(message, 0)
	assert.Equal(t, datatypes.MissingStartTime, err, "Check for missing start time")

	message = []byte("{\"stoneID\":\"test1234\",\"startTime\":1548258873000}")
	_, err = parseFlag2(message, 0)
	assert.Equal(t, datatypes.MissingEndTime, err, "Check for missing end time")
}