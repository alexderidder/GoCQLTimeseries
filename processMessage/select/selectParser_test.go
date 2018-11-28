package _select

import (
	"GoCQLTimeSeries/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelectFlagParser(t *testing.T) {
	message := &[]byte{1, 0, 0, 0}
	_, err := Parse(message)
	assert.Equal(t, model.Error{301 , "unexpected end of JSON input"}, err, "Parseflag returns Marshall error, so flag is known")
	message = &[]byte{0, 0, 0, 0}
	_, err = Parse(message)
	assert.Equal(t, model.Error{150 , "Flag doesn't exists"}, err, "Parseflag returns Flag doesn't exists")
	message = &[]byte{0, 0, 0, 0}
	_, err = Parse(message)
	assert.Equal(t, model.FlagNoExist, err, "Parseflag returns Flag doesn't exists")
	message = &[]byte{}
	_, err = Parse(message)
	assert.Equal(t, model.MessageNoLengthForFlag, err, "Parseflag returns message no length for flag")
}