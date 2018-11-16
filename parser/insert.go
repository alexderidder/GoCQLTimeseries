package parser

import (
	"../database"
	"../model"
	"encoding/binary"
	"encoding/json"
)

// I'd suggest naming this a InsertQuery, or InsertMessage etc.. Just Insert seems like a function call to insert something.

// THE RESPONSIBILITY OF IMPLEMENTING THE Command INTERFACE LIES WITH THESE FILES, NOT FORCING THEM INTO A METHOD EXPECTING A COMMAND!
// THE RESPONSIBILITY OF IMPLEMENTING THE FlagMethods INTERFACE LIES WITH THESE FILES, NOT FORCING THEM INTO A METHOD EXPECTING A COMMAND!
type Insert struct {
	flag    []byte
	message []byte
	request *model.InsertJSON
}

//!! InsertFlag1 is not a good variable name. You know this!
type InsertFlag1 struct {
	insert *Insert
}

// i is not a good variable name. You know this!
func (i Insert) parseFlag() []byte {
	flag := model.ByteToUint32(i.flag, 0)
	switch flag {
	case 1:
		// what is happening here!?
		return i.executeMethodsPerFlag(InsertFlag1{&i}) // why don't we put "flag" into this?
	default:
		return model.FlagNoExist.MarshallErrorAndAddFlag()
	}
}

// i is not a good variable name. You know this!
// test2 is not a good variable name. You know this!
func (i Insert) executeMethodsPerFlag(test2 FlagMethods) []byte {
	// this is inconsistent with the checkParameters, pick one and stick to it.
	if !test2.marshalBytes() {
		return model.ErrorMarshal.MarshallErrorAndAddFlag()
	}
	error := test2.checkParameters()
	if !error.IsNull() {
		return error.MarshallErrorAndAddFlag()
	}
	return test2.databaseInteraction()
}

// i is not a good variable name. You know this!
func (i InsertFlag1) marshalBytes() bool {
	err := json.Unmarshal(i.insert.message, &i.insert.request)
	if err != nil {
		return false
	}
	return true
}

// i is not a good variable name. You know this!
func (i InsertFlag1) checkParameters() model.Error {

	if !i.insert.request.StoneID.Valid {
		return model.MissingStoneID
	}

	if len(i.insert.request.Data) == 0 {
		return model.MissingData
	}

	return model.NoError
}

// i is not a good variable name. You know this!
func (i InsertFlag1) databaseInteraction() []byte {
	err := database.Insert(i.insert.request)
	if !err.IsNull() {
		return err.MarshallErrorAndAddFlag()
	}

	// code like this is used a lot. Try having a util method like:
	// return Util.uint32_to_byteArray(2) or something.
	resultCode := make([]byte, 4) // magic number,
	binary.LittleEndian.PutUint32(resultCode, 2)
	return resultCode
}
