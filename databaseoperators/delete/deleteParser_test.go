package delete

import (
	"GoCQLTimeSeries/datatypes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleteFlagParser(t *testing.T) {
	message := []byte{1, 0, 0, 0}
	_, err := Parse(message)
	assert.Equal(t, datatypes.Error{301 , "unexpected end of JSON input"}, err, "Parseflag returns Marshall error, so flag is known")
	message = []byte{0, 0, 0, 0}
	_, err = Parse(message)
	assert.Equal(t, datatypes.FlagNoExist, err, "Parseflag returns Flag doesn't exists")
	message = []byte{}
	_, err = Parse(message)
	assert.Equal(t, datatypes.MessageNoLengthForFlag, err, "Parseflag returns message no length for flag")
}
