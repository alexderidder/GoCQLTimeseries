package parser

import (
	"../database"
	"../model"
	"encoding/binary"
	"encoding/json"
)

// I'd suggest naming this a DeleteQuery, or DeleteMessage etc..

// THE RESPONSIBILITY OF IMPLEMENTING THE Command INTERFACE LIES WITH THESE FILES, NOT FORCING THEM INTO A METHOD EXPECTING A COMMAND!
// THE RESPONSIBILITY OF IMPLEMENTING THE FlagMethods INTERFACE LIES WITH THESE FILES, NOT FORCING THEM INTO A METHOD EXPECTING A COMMAND!
type Delete struct {
	flag    []byte
	message []byte
	request *model.DeleteJSON
}

// check variable names! See Insert for suggestions on what's wrong here.
type DeleteFlag1 struct {
	delete *Delete
}

func (d Delete) parseFlag() []byte {
	flag := model.ByteToUint32(d.flag, 0)
	switch flag {
	case 1:
		return d.executeMethodsPerFlag(DeleteFlag1{&d})
	default:
		return model.FlagNoExist.MarshallErrorAndAddFlag()
	}
}

func (d Delete) executeMethodsPerFlag(test2 FlagMethods) []byte {

	if !test2.marshalBytes() {
		return model.ErrorMarshal.MarshallErrorAndAddFlag()
	}
	error := test2.checkParameters()
	if !error.IsNull() {
		return error.MarshallErrorAndAddFlag()
	}
	return test2.databaseInteraction()
}

func (d DeleteFlag1) marshalBytes() bool {
	err := json.Unmarshal(d.delete.message, &d.delete.request)
	if err != nil {
		return false
	}
	return true
}

func (d DeleteFlag1) checkParameters() model.Error {

	if !d.delete.request.StoneID.Valid {
		return model.MissingStoneID
	}

	if len(d.delete.request.Types) == 0 {
		return model.MissingType
	}
	d.delete.request.Types = checkUnknownAndDuplicatedTypes(d.delete.request.Types)

	return model.NoError
}

func (d DeleteFlag1) databaseInteraction() []byte {
	err := database.Delete(d.delete.request)
	if !err.IsNull() {
		return err.MarshallErrorAndAddFlag()
	}
	resultCode := make([]byte, 4)
	binary.LittleEndian.PutUint32(resultCode, 2)
	return resultCode
}
