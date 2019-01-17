package insert

import (
	"GoCQLTimeSeries/model"
	"GoCQLTimeSeries/util"
)

func Parse(message []byte) (model.Execute, model.Error) {
	if len(message) < 4 {
		return nil, model.MessageNoLengthForFlag
	}
	flag := util.GetUInt32FromIndex(0, message)
	indexOfMessage := 4
	switch flag {
	case 1:
		return parseFlag1(message, indexOfMessage)
	case 2:
		return parseFlag2(message, indexOfMessage)
	default:
		return nil, model.FlagNoExist
	}
}