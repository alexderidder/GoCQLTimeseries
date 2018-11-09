package parser

import (
	"CrownstoneServer/database"
	"CrownstoneServer/model"
	"encoding/binary"
	"encoding/json"
)



type Delete struct {
	message    []byte
	request *model.DeleteJSON
}


func (d Delete) parseFlag() bool {
	flag := model.ByteToInt(d.message, 0)
	switch flag {
	case 1:
		return true
	default:
		return false
	}
}
func (d Delete) parseJSON() bool {
	err := json.Unmarshal(d.message[4:], &d.request)
	if err != nil {
		return false
	}
	return true
}
func (d Delete) checkParameters() model.Error {

	if !d.request.StoneID.Valid {
		return model.Error{100, "StoneID is missing"}
	}

	if len(d.request.Types) == 0 {
		return model.Error{100, "Type is missing"}
	}
	d.request.Types = checkUnknownAndDuplicatedTypes(d.request.Types)

	return model.Null
}

func (d Delete) parseJSONToDatabaseQueries() []byte {
	err := database.Delete(d.request)
	if !err.IsNull() {
		return err.MarshallErrorAndAddFlag()
	}
	resultCode := make([]byte, 4)
	binary.LittleEndian.PutUint32(resultCode, 2)
	return resultCode
}


