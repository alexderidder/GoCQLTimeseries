package datatypes

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckError(t *testing.T) {

	error := Error{0, "Hallo"}
	assert.True(t, error.IsNull(),    "Check error is null")

	error = Error{1, "Hallo"}
	assert.False(t, error.IsNull(),    "Check error is not null")

}
func TestMarshallErrorFlag(t *testing.T) {

	error := Error{1, "Hallo"}
	bytes := error.marshalError()
	assert.Equal(t, "{\"code\":1,\"message\":\"Hallo\"}", string(bytes),    "Check empty error is marshalled correctly")

	error = Error{}
	bytes = error.marshalError()
	assert.Equal(t, "{\"code\":0,\"message\":\"\"}", string(bytes),    "Check empty error is marshalled correctly")

}

func TestMarshallErrorAndAddFlag(t *testing.T) {

	error := Error{1, "Hallo"}
	bytes := error.MarshallErrorAndAddFlag()
	assert.Equal(t, "d\x00\x00\x00{\"code\":1,\"message\":\"Hallo\"}", string(bytes),    "Check empty error is marshalled correctly")

	error = Error{}
	bytes = error.MarshallErrorAndAddFlag()
	assert.Equal(t, "d\x00\x00\x00{\"code\":0,\"message\":\"\"}", string(bytes),    "Check empty error is marshalled correctly")

}
