package util

import (
	"strings"
)

const (
	UnitW      string = "w"
	Unitpf     string = "pf"
	UnitWAndpf string = "w_pf"
	UnitkWh    string = "kwh"
)


func CheckUnknownAndDuplicatedTypes(request []string) []string {
	var typeList = []bool{false, false, false}
	for _, v := range request {
		switch strings.ToLower(v) {
		case UnitW:
			typeList[0] = true
		case Unitpf:
			typeList[1] = true
		case UnitkWh:
			typeList[2] = true
		}
	}
	typePerQuery := make([]string, 2)
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

	return typePerQuery

}

func ReturnAndRemoveUint32FromByteArrayByIndex(index uint32, message *[]byte) uint32{
	flag := ByteToUint32(*message, int(index))
	tempMessage := *message
	*message = tempMessage[4:]
	return flag
}