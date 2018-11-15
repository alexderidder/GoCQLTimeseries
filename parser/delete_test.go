package parser

import (
	"github.com/alexderidder/GoCQLTimeseries/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

//func TestDeleteFlagParser(t *testing.T) {
//	test := model.DeleteJSON{}
//	delete := Delete{[]byte{1, 0, 0, 0}, &test}
//	assert.True(t, delete.parseFlag(), "Parseflag returns true when flag is 1")
//	delete.message = []byte{0, 0, 0, 0}
//	assert.False(t, delete.parseFlag(), "Parseflag returns false when flag is 0")
//	delete.message = []byte{}
//	assert.False(t, delete.parseFlag(), "Parseflag returns false when flag is not given")
//}

func TestDeleteParseJSON(t *testing.T) {
	test := model.DeleteJSON{}
	delete := Delete{[]byte{1, 0, 0, 0}, []byte("{\"stoneID\":\"bf82e78d-24a2-470d-abb8-9e0a2720619f\",\"types\":[\"w\",\"pf\",\"kwh\"],\"startTime\":\"1995-03-01T00:00:00Z\",\"endTime\":\"1995-02-01T00:00:00Z\"}"), &test}
	flag := DeleteFlag1{&delete}
	assert.Equal(t, true, flag.marshalBytes(), "Marshal byte went wrong")
	assert.Equal(t, "bf82e78d-24a2-470d-abb8-9e0a2720619f", test.StoneID.Value.String(), "StoneID is not bf82e78d-24a2-470d-abb8-9e0a2720619f")
	assert.Equal(t, []string{"w", "pf", "kwh"}, test.Types, "Types is [w, pf, kwh] ")
	assert.Equal(t, "1995-03-01 00:00:00 +0000 UTC", test.StartTime.String(), "Starttime isnt 1995-03-01T00:00:00Z ")
	assert.Equal(t, "1995-02-01 00:00:00 +0000 UTC", test.EndTime.String(), "Endtime isnt 1995-03-01T00:00:00Z ")
}

func TestDeleteRequirements(t *testing.T) {
	test := model.DeleteJSON{}
	delete := Delete{[]byte{1, 0, 0, 0}, []byte("{\"stoneID\":\"bf82e78d-24a2-470d-abb8-9e0a2720619f\",\"types\":[\"w\",\"pf\",\"kwh\"],\"startTime\":\"1995-03-01T00:00:00Z\",\"endTime\":\"1995-02-01T00:00:00Z\"}"), &test}
	flag := DeleteFlag1{&delete}
	assert.Equal(t, true, flag.marshalBytes(), "Marshal byte went wrong")
	assert.Equal(t, model.NoError, flag.checkParameters(), "Check if no error")

	test = model.DeleteJSON{}
	delete = Delete{[]byte{1, 0, 0, 0}, []byte("{\"types\":[\"w\",\"pf\",\"kwh\"],\"startTime\":\"1995-03-01T00:00:00Z\",\"endTime\":\"1995-02-01T00:00:00Z\"}"), &test}
	assert.Equal(t, true, flag.marshalBytes(), "Marshal byte went wrong")
	assert.Equal(t, model.MissingStoneID, flag.checkParameters(), "Check if stone ID is missing")

	test = model.DeleteJSON{}
	delete = Delete{[]byte{1, 0, 0, 0}, []byte("{\"stoneID\":\"bf82e78d-24a2-470d-abb8-9e0a2720619f\",\"startTime\":\"1995-03-01T00:00:00Z\",\"endTime\":\"1995-02-01T00:00:00Z\"}"), &test}
	assert.Equal(t, true, flag.marshalBytes(), "Marshal byte went wrong")
	assert.Equal(t, model.MissingType, flag.checkParameters(), "Check if types are missing")
}
