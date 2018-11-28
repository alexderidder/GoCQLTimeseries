package insert

import (
	"GoCQLTimeSeries/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInsertFlag1JSON(t *testing.T) {
	message := []byte("{\"stoneID\":\"bf82e78d-24a2-470d-abb8-9e0a2720619f\",\"data\":[{\"time\":\"2018-11-14T10:16:26.7220151+01:00\",\"kWh\":3.3228004,\"watt\":3.0233014,\"pf\":4.702545}]}")
	_, err := parseFlag1(&message)
	assert.Equal(t, model.NoError, err, "Check if no Error")
	message = []byte("")
	_, err = parseFlag1(&message)
	assert.Equal(t, model.Error{301, "unexpected end of JSON input"}, err, "Check for marshal error unexpected end")

	message =  []byte("{\"data\":[{\"time\":\"2018-11-14T10:16:26.7220151+01:00\",\"kWh\":3.3228004,\"watt\":3.0233014,\"pf\":4.702545}]}")
	_, err = parseFlag1(&message)
	assert.Equal(t, model.MissingStoneID, err, "Check for missing stone id")
	message =  []byte("{\"stoneID\":\"bf82e78d-24a2-470d-abb8-9e0a2720619f\"}")
	_, err = parseFlag1(&message)
	assert.Equal(t, model.MissingData, err, "Check if data is  missing")
}