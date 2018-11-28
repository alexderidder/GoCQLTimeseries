package _select

import (
	"GoCQLTimeSeries/model"
	"GoCQLTimeSeries/util"
)
const(
	KWH = 1
	Watt_PowerFactor = 2
)
func Parse(message *[]byte) (model.Execute, model.Error) {
	if len(*message) < 4 {
		return nil, model.MessageNoLengthForFlag
	}
	flag := util.ReturnAndRemoveUint32FromByteArrayByIndex(0, message)
	switch flag {
	case KWH:
		return parseFlag1(message)
	case Watt_PowerFactor:
		return parseFlag2(message)
	default:
		return nil, model.FlagNoExist
	}
}