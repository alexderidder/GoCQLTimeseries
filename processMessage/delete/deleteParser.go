package delete

import (
	"GoCQLTimeSeries/model"
	"GoCQLTimeSeries/util"
)

func  Parse(message *[]byte) (model.Execute, model.Error) {
	if len(*message) < 4 {
		return nil, model.MessageNoLengthForFlag
	}
	flag := util.ReturnAndRemoveUint32FromByteArrayByIndex(0, message)
	switch flag {
	case 1:
		return parseFlag1(message)
	default:
		return nil, model.FlagNoExist
	}
}
