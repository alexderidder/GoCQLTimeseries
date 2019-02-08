package databaseoperators

import (
	"GoCQLTimeSeries/datatypes"
	"GoCQLTimeSeries/databaseoperators/delete"
	"GoCQLTimeSeries/databaseoperators/insert"
	"GoCQLTimeSeries/databaseoperators/select"
	"GoCQLTimeSeries/model"
)

func ParseOpCode(opCode uint32, message []byte) (model.Execute, datatypes.Error) {
	switch opCode {
	case 100:
		executeObject, err := insert.Parse(message)
		return executeObject, err
	case 200:
		executeObject, err := _select.Parse(message)
		return executeObject, err
	case 500:
		executeObject, err := delete.Parse(message)
		return executeObject, err
	default:
		return nil, datatypes.Error{10, "Server doesn't recognise opcode"}
	}
}
