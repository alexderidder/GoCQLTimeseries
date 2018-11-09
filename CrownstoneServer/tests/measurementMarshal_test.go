package tests

import (
	"CrownstoneServer/parser"
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateUser(t *testing.T){
	value := MarshalMeasurement(parser.Data{time.Time{}, 3.3 })
	assert.True(t, bytes.Equal(value,  []byte("{\"time\":\"0001-01-01T00:00:00Z\",\"value\":3.3}")), "Check if equal")

}
