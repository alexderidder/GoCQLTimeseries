package parser

import (
	"CrownstoneServer/database"
	"CrownstoneServer/model"
	"encoding/binary"
	"encoding/json"
)

type Insert struct {
	message []byte
	request *model.InsertJSON
}

func (i Insert) parseFlag() bool {
	flag := model.ByteToUint32(i.message, 0)
	switch flag {
	case 1:
		i.message = i.message[4:]
		return true
	default:
		return false
	}
}
func (i Insert) parseJSON() bool {
	err := json.Unmarshal(i.message, &i.request)
	if err != nil {
		return false
	}
	return true
}
func (i Insert) checkParameters() model.Error {

	if !i.request.StoneID.Valid {
		return model.MissingStoneID
	}

	if len(i.request.Data) == 0 {
		return model.MissingData
	}

	return model.Null
}

func (i Insert) parseJSONToDatabaseQueries() []byte {
	err := database.Insert(i.request)
	if !err.IsNull() {
		return err.MarshallErrorAndAddFlag()
	}
	resultCode := make([]byte, 4)
	binary.LittleEndian.PutUint32(resultCode, 1)
	return resultCode
}
