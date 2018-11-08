package parser

import (
	"CrownstoneServer/server/database/connector"
	"encoding/binary"
	"encoding/json"
	"strings"
)

type Error struct {
	Code    uint32 `json:"code"`
	Message string `json:"message"`
}

type Command interface {
	parseFlag() []byte
}

func ParseOpCode(opCode uint32, message []byte) []byte {
	if cassandra.Session == nil{
		return ParseError(300, "There is no connection between GoCQL-driver and Cassandra")
	}
	switch opCode {
	//TODO: insert
	case 100:
		i := insert{message}
		return parser(i)
	case 200:
		s := selectCommand{message}
		return parser(s)
	case 500:
		d := delete{message}
		return  parser(d)
		//TODO: Research delete management
	default:
		return ParseError(10, "Server doesn't recognise opcode")
	}
}

func parser(c Command) []byte {
	return c.parseFlag()
}

func ParseHeader(request []byte) (uint32, uint32, uint32, uint32) {
	result := make([]uint32, 4)
	for i := 0; i < 4; i++ {
		result[i] = byteToInt(request, i*4)
	}

	return result[0], result[1], result[2], result[3]
}



func byteToInt(request []byte, beginIndex int) uint32 {
	var result uint32
	result |= uint32(request[beginIndex])
	beginIndex++
	result |= uint32(request[beginIndex]) << 8
	beginIndex++
	result |= uint32(request[beginIndex]) << 16
	beginIndex++
	result |= uint32(request[beginIndex]) << 24
	return result
}

func ParseError(code uint32, message string) []byte {
	var error Error
	if code == 0 && len(message) == 0 {
		error = Error{500, "Server created error without code and message"}
	} else if code == 0 {
		error = Error{500, "Server created error without code"}
	} else if len(message) == 0 {
		error = Error{500, "Server created error without message"}
	} else {
		error = Error{code, message}
	}

	errBytes, _ := json.Marshal(error)
	//TODO: Marshal error
	errCode := make([]byte, 4)
	binary.LittleEndian.PutUint32(errCode, 100)
	return append(errCode, errBytes...)
}


func checkUnknownAndDuplicatedTypes(request []string) (bool, [2]string) {
	var typeList = []bool{false, false, false}
	for _, v := range request {
		switch strings.ToLower(v) {
		case UnitW:
			typeList[0] = true
		case Unitpf:
			typeList[1] = true
		case UnitkWh:
			typeList[2] = true
		default:
			return false, [2]string{}
		}

	}
	var typePerQuery [2]string
	if typeList[0] && typeList[1] {
		typePerQuery[0] = UnitWAndpf
	} else if typeList[0] {
		typePerQuery[0] = UnitW
	} else if typeList[1] {
		typePerQuery[0] = Unitpf
	}

	if typeList[2] {
		typePerQuery[1] = UnitkWh
	}

	return true, typePerQuery

}

