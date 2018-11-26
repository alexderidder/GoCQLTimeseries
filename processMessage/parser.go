package processMessage

import (
	"../model"
	"../util"
	"./delete"
	"./insert"
	"./select"
)

type Operation interface {
	Parse(*[]byte) (Execute, model.Error)
}

type Execute interface {
	Execute() ([]byte, model.Error)
}

//@Alex is this what you expected? I made for every Operator + flag a object.
func ParseOpCode(opCode uint32, message *[]byte) (Execute, model.Error){
	switch opCode {
	case 100:
		executeObject, err := parseInsert(message)
		return executeObject, err
	case 200:
		executeObject, err := parseSelect(message)
		return executeObject, err
	case 500:
		executeObject, err := parseDelete(message)
		return executeObject, err
	default:
		return nil, model.Error{10, "Server doesn't recognise opcode"}
	}
}

func  parseInsert(message *[]byte) (Execute, model.Error) {
	flag := readFlagAndRemoveFromMessage(message)
	switch flag {
	case 1:
		return insert.ParseFlag1(message)
	default:
		return nil, model.FlagNoExist
	}
}

func  parseDelete(message *[]byte) (Execute, model.Error) {
	flag := readFlagAndRemoveFromMessage(message)
	switch flag {
	case 1:
		return delete.ParseFlag1(message)
	default:
		return nil, model.FlagNoExist
	}
}

func  parseSelect(message *[]byte) (Execute, model.Error) {
	flag := readFlagAndRemoveFromMessage(message)
	switch flag {
	case 1:
		return _select.ParseFlag1(message)
	default:
		return nil, model.FlagNoExist
	}
}

func readFlagAndRemoveFromMessage(message *[]byte) uint32{
	flag := util.ByteToUint32(*message, 0)
	tempMessage := *message
	*message = tempMessage[4:]
	return flag
}
