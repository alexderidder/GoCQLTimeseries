package processMessage

import (
	"GoCQLTimeSeries/model"
	"GoCQLTimeSeries/processMessage/delete"
	"GoCQLTimeSeries/processMessage/insert"
	"GoCQLTimeSeries/processMessage/select"
)

func ParseOpCode(opCode uint32, message *[]byte) (model.Execute, model.Error) {
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
		return nil, model.Error{10, "Server doesn't recognise opcode"}
	}
}
