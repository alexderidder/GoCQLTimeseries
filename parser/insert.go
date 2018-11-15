package parser

import (
	"github.com/alexderidder/GoCQLTimeseries/database"
	"github.com/alexderidder/GoCQLTimeseries/model"
	"encoding/binary"
	"encoding/json"
)

type Insert struct {
	flag []byte
	message []byte
	request model.InsertJSON
}

type InsertFlag1 struct {
	insert *Insert
}


func (i Insert) parseFlag() []byte {
	flag := model.ByteToUint32(i.flag, 0)
	switch flag {
	case 1:
		return i.executeMethodsPerFlag(InsertFlag1{&i})
	default:
		return model.FlagNoExist.MarshallErrorAndAddFlag()
	}
}

func (i Insert) executeMethodsPerFlag(test2 FlagMethods) []byte {

	if !test2.marshalBytes() {
		return model.ErrorMarshal.MarshallErrorAndAddFlag()
	}
	error := test2.checkParameters()
	if !error.IsNull(){
		return error.MarshallErrorAndAddFlag()
	}
	return test2.databaseInteraction()
}

func (i InsertFlag1) marshalBytes() bool {
	err := json.Unmarshal(i.insert.message, &i.insert.request)
	if err != nil {
		return false
	}
	return true
}
func (i InsertFlag1) checkParameters() model.Error {

	if !i.insert.request.StoneID.Valid {
		return model.MissingStoneID
	}

	if len(i.insert.request.Data) == 0 {
		return model.MissingData
	}

	return model.NoError
}

func (i InsertFlag1) databaseInteraction() []byte {
	err := database.Insert(i.insert.request)
	if !err.IsNull() {
		return err.MarshallErrorAndAddFlag()
	}
	resultCode := make([]byte, 4)
	binary.LittleEndian.PutUint32(resultCode, 1)
	return resultCode
}
