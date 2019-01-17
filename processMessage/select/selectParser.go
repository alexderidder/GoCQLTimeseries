package _select

import (
	"GoCQLTimeSeries/model"
	"GoCQLTimeSeries/util"
)
const(
	KWH = 1
	Watt_PowerFactor = 2
)
func Parse(message []byte) (model.Execute, model.Error) {
	if len(message) < 4 {
		return nil, model.MessageNoLengthForFlag
	}
	flag := util.GetUInt32FromIndex(0, message)
	indexOfMessage := 4
	switch flag {
	case KWH:
		return parseFlag1(message, indexOfMessage)
	case Watt_PowerFactor:
		return parseFlag2(message, indexOfMessage)
	default:
		return nil, model.FlagNoExist
	}
}