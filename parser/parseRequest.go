package parser

import (
	"../model"
	"strings"
)

type Command interface {
	parseFlag() []byte
	executeMethodsPerFlag(FlagMethods) []byte
}

type FlagMethods interface {
	marshalBytes() bool
	checkParameters() model.Error
	databaseInteraction() []byte
}

// why does this method not expect a flag and a message? You seem to hardcode the first 4 bytes of the message as flag anyway
// this method should not have to know this.
// Alternatively, this method should only know about the opcodes, not the content of the message. In that case, the Insert struct should know this.
func ProcessOpCodeAndReceivedMessage(opCode uint32, message []byte) []byte {
	switch opCode {
	//TODO: insert
	case 100:
		// i, s and d are not really good variable names.
		i := Insert{message[:4], message[4:], &model.InsertJSON{}}
		return parser(i)
	case 200:
		s := Get{message[:4], message[4:], &model.RequestSelectJSON{}}
		return parser(s)
	case 500:
		d := Delete{message[:4], message[4:], &model.DeleteJSON{}}
		return parser(d)
		//TODO: Research delete management
	default:
		return model.Error{10, "Server doesn't recognise opcode"}.MarshallErrorAndAddFlag()
	}
}

// what is the purpose of this function?
func parser(c Command) []byte {
	// all we do here is "parseFlag", while the method that's calling this says process message.
	// After parsing, the message content is executed. This is an unclear side effect of calling this parse method.
	return c.parseFlag()
}

func checkUnknownAndDuplicatedTypes(request []string) []string {
	var typeList = []bool{false, false, false}
	for _, v := range request {
		switch strings.ToLower(v) {
		case model.UnitW:
			typeList[0] = true
		case model.Unitpf:
			typeList[1] = true
		case model.UnitkWh:
			typeList[2] = true
		}
	}
	typePerQuery := make([]string, 2)
	if typeList[0] && typeList[1] {
		typePerQuery[0] = model.UnitWAndpf
	} else if typeList[0] {
		typePerQuery[0] = model.UnitW
	} else if typeList[1] {
		typePerQuery[0] = model.Unitpf
	}

	if typeList[2] {
		typePerQuery[1] = model.UnitkWh
	}

	return typePerQuery

}
