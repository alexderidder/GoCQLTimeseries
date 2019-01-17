package delete

import (
	"GoCQLTimeSeries/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleteFlag2JSON(t *testing.T) {
	message := []byte("{\"stoneID\":\"bf82e78d-24a2-470d-abb8-9e0a2720619f\",\"types\":[\"w\",\"pf\",\"kwh\"],\"startTime\":\"1995-03-01T00:00:00Z\",\"endTime\":\"1995-02-01T00:00:00Z\"}")
	_, err := parseFlag2(message, 0)
	assert.Equal(t, model.NoError, err, "Check if no Error")
	message = []byte("")
	_, err = parseFlag2(message, 0)
	assert.Equal(t, model.Error{301, "unexpected end of JSON input"}, err, "Check for marshal error unexpected end")

	message = []byte("{\"types\":[\"w\",\"pf\",\"kwh\"],\"startTime\":\"1995-03-01T00:00:00Z\",\"endTime\":\"1995-02-01T00:00:00Z\"}")
	_, err = parseFlag2(message, 0)
	assert.Equal(t, model.MissingStoneID, err, "Check for missing stone id")

}
