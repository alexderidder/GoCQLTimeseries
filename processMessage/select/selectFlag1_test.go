package _select

import (
	"GoCQLTimeSeries/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelectFlag1JSON(t *testing.T) {
	message := []byte("{\"stoneIDs\":[\"bf82e78d-24a2-470d-abb8-9e0a2720619f\"],\"types\":[\"w\",\"pf\",\"kwh\"],\"startTime\":\"2018-11-14T10:40:17.5485474+01:00\",\"endTime\":\"2018-11-14T11:15:17.5485474+01:00\",\"interval\":0}")
	_, err := parseFlag1(&message)
	assert.Equal(t, model.NoError, err, "Check if no Error")
	message = []byte("")
	_, err = parseFlag1(&message)
	assert.Equal(t, model.Error{301, "unexpected end of JSON input"}, err, "Check for marshal error unexpected end")

	message =  []byte("{\"types\":[\"w\",\"pf\",\"kwh\"],\"startTime\":\"2018-11-14T10:40:17.5485474+01:00\",\"endTime\":\"2018-11-14T11:15:17.5485474+01:00\",\"interval\":0}")
	_, err = parseFlag1(&message)
	assert.Equal(t, model.MissingStoneID, err, "Check for missing stone id")


}