package parser

import (
	"CrownstoneServer/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInsertFlagParser(t *testing.T) {
	test := model.InsertJSON{}
	insert := Insert{[]byte{1, 0, 0, 0}, &test}
	assert.True(t, insert.parseFlag(), "Parseflag returns true when flag is 1")
	insert.message = []byte{0, 0, 0, 0}
	assert.False(t, insert.parseFlag(), "Parseflag returns false when flag is 0")
	insert.message = []byte{}
	assert.False(t, insert.parseFlag(), "Parseflag returns false when flag is not given")
}

func TestInsertParseJSON(t *testing.T) {
	test := model.InsertJSON{}
	insert := Insert{[]byte("{\"stoneID\":\"bf82e78d-24a2-470d-abb8-9e0a2720619f\",\"data\":[{\"time\":\"2018-11-14T10:16:26.7220151+01:00\",\"kWh\":3.3228004,\"watt\":3.0233014,\"pf\":4.702545}]}"), &test}
	assert.Equal(t, true, insert.parseJSON(), "Marshal byte went wrong")
	assert.Equal(t, "bf82e78d-24a2-470d-abb8-9e0a2720619f", test.StoneID.Value.String(), "StoneID is not bf82e78d-24a2-470d-abb8-9e0a2720619f")
	assert.Equal(t, 1, len(test.Data), "Data hasn't length 1 ")
	assert.Equal(t, "2018-11-14 10:16:26.7220151 +0100 CET", test.Data[0].Time.String(), "Time isnt 2018-11-14 10:16:26.7220151 +0100 CET")
	assert.Equal(t, float32(3.3228004), test.Data[0].KWH.Value, "kwh = 3.3228004")
}

func TestInsertRequirements(t *testing.T) {
	test := model.InsertJSON{}
	insert := Insert{[]byte("{\"stoneID\":\"bf82e78d-24a2-470d-abb8-9e0a2720619f\",\"data\":[{\"time\":\"2018-11-14T10:16:26.7220151+01:00\",\"kWh\":3.3228004,\"watt\":3.0233014,\"pf\":4.702545}]}"), &test}
	assert.Equal(t, true, insert.parseJSON(), "Marshal byte went wrong")
	assert.Equal(t, model.Null, insert.checkParameters(), "Check if no error")

	test = model.InsertJSON{}
	insert = Insert{[]byte("{\"data\":[{\"time\":\"2018-11-14T10:16:26.7220151+01:00\",\"kWh\":3.3228004,\"watt\":3.0233014,\"pf\":4.702545}]}"), &test}
	assert.Equal(t, true, insert.parseJSON(), "Marshal byte went wrong")
	assert.Equal(t, model.MissingStoneID, insert.checkParameters(), "Check if stone ID is missing")

	test = model.InsertJSON{}
	insert = Insert{[]byte("{\"stoneID\":\"bf82e78d-24a2-470d-abb8-9e0a2720619f\"}"), &test}
	assert.Equal(t, true, insert.parseJSON(), "Marshal byte went wrong")
	assert.Equal(t, model.MissingData, insert.checkParameters(), "Check if data is missing")
}
