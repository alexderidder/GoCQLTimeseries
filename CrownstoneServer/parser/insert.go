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
	flag := model.ByteToInt(i.message, 0)
	switch flag {
	case 1:
		return true
	default:
		return false
	}
}
func (i Insert) parseJSON() bool {
	err := json.Unmarshal(i.message[4:], &i.request)
	if err != nil {
		return false
	}
	return true
}
func (i Insert) checkParameters() model.Error {

	if !i.request.StoneID.Valid {
		return model.Error{100, "StoneID is missing"}
	}

	if len(i.request.Data) == 0 {
		return model.Error{100, "Data is missing"}
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
