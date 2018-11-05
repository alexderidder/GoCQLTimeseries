package tests

import (
	"CrownstoneServer/parser"
	"encoding/json"
)


func MarshalMeasurement(value parser.Data) []byte {
	result, _ := json.Marshal(value)
	return result
}
