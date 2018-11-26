package _select

//import (
//	"CrownstoneServer/model"
//	"github.com/stretchr/testify/assert"
//	"testing"
//)
//
////func TestSelectFlagParser(t *testing.T) {
////	test := model.RequestJSON{}
////	insert := Get{[]byte{1, 0, 0, 0}, &test}
////	flag1 := GetFlag1{&insert}
////	assert.True(t, flag1.parseFlag(), "Parseflag returns true when flag is 1")
////	insert.message = []byte{0, 0, 0, 0}
////	assert.False(t, insert.parseFlag(), "Parseflag returns false when flag is 0")
////	insert.message = []byte{}
////	assert.False(t, insert.parseFlag(), "Parseflag returns false when flag is not given")
////}
//
//func TestSelectParseJSON(t *testing.T) {
//	test := model.RequestJSON{}
//	get := Get{[]byte{1, 0, 0, 0}, []byte("{\"stoneIDs\":[\"bf82e78d-24a2-470d-abb8-9e0a2720619f\"],\"types\":[\"w\",\"pf\",\"kwh\"],\"startTime\":\"2018-11-14T10:40:17.5485474+01:00\",\"endTime\":\"2018-11-14T11:15:17.5485474+01:00\",\"interval\":0}"), &test}
//	flag1 := GetFlag1{&get}
//	assert.Equal(t, true, flag1.marshalBytes(), "Marshal byte went wrong")
//	assert.Equal(t, "bf82e78d-24a2-470d-abb8-9e0a2720619f", test.StoneIDs[0].Value.String(), "StoneID is not bf82e78d-24a2-470d-abb8-9e0a2720619f")
//	assert.Equal(t, []string{"w", "pf", "kwh"}, test.Types, "Types is [w, pf, kwh] ")
//	assert.Equal(t, "2018-11-14 10:40:17.5485474 +0100 CET", test.StartTime.String(), "Time isn 2018-11-14 10:40:17.5485474 +0100 CET")
//	assert.Equal(t, "2018-11-14 11:15:17.5485474 +0100 CET", test.EndTime.String(), "Time isn 2018-11-14 11:15:17.5485474 +0100 CET")
//}
//
//func TestSelectRequirements(t *testing.T) {
//	test := model.RequestJSON{}
//	get := Get{[]byte{1, 0, 0, 0}, []byte("{\"stoneIDs\":[\"bf82e78d-24a2-470d-abb8-9e0a2720619f\"],\"types\":[\"w\",\"pf\",\"kwh\"],\"startTime\":\"2018-11-14T10:40:17.5485474+01:00\",\"endTime\":\"2018-11-14T11:15:17.5485474+01:00\",\"interval\":0}"), &test}
//	flag1 := GetFlag1{&get}
//	assert.Equal(t, true, flag1.marshalBytes(), "Marshal byte went wrong")
//	assert.Equal(t, model.NoError, flag1.checkParameters(), "Check if no error")
//
//	test = model.RequestJSON{}
//	get = Get{[]byte{1, 0, 0, 0}, []byte("{\"types\":[\"w\",\"pf\",\"kwh\"],\"startTime\":\"2018-11-14T10:40:17.5485474+01:00\",\"endTime\":\"2018-11-14T11:15:17.5485474+01:00\",\"interval\":0}"), &test}
//	assert.Equal(t, true, flag1.marshalBytes(), "Marshal byte went wrong")
//	assert.Equal(t, model.MissingStoneID, flag1.checkParameters(), "Check if stone ID is missing")
//
//	test = model.RequestJSON{}
//	get = Get{[]byte{1, 0, 0, 0}, []byte("{\"stoneIDs\":[\"bf82e78d-24a2-470d-abb8-9e0a2720619f\"],\"startTime\":\"2018-11-14T10:40:17.5485474+01:00\",\"endTime\":\"2018-11-14T11:15:17.5485474+01:00\",\"interval\":0}"), &test}
//
//	assert.Equal(t, true, flag1.marshalBytes(), "Marshal byte went wrong")
//	assert.Equal(t, model.MissingType, flag1.checkParameters(), "Check if data is missing")
//}
