package _select

import (
	"GoCQLTimeSeries/datatypes"
	"GoCQLTimeSeries/model"
	"GoCQLTimeSeries/util"
)
const(
	KWH = 1
	Watt_PowerFactor = 2
)
func Parse(message []byte) (model.Execute, datatypes.Error) {
	if len(message) < 4 {
		return nil, datatypes.MessageNoLengthForFlag
	}
	flag := util.GetUInt32FromIndex(0, message)
	indexOfMessage := 4
	switch flag {
	case KWH:
		return parseFlag1(message, indexOfMessage)
	case Watt_PowerFactor:
		return parseFlag2(message, indexOfMessage)
	default:
		return nil, datatypes.FlagNoExist
	}
}