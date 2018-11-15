package parser

import (
	"../database"
	"../model"
	"encoding/binary"
	"encoding/json"
)

const (
	UnitW      string = "w"
	Unitpf     string = "pf"
	UnitWAndpf string = "w_pf"
	UnitkWh    string = "kwh"
)

type Get struct {
	flag []byte
	message []byte
	request *model.RequestSelectJSON
}

type GetFlag1 struct {
	get *Get
}

func (g Get) parseFlag() []byte {
	flag := model.ByteToUint32(g.flag, 0)
	switch flag {
	case 1:
		return g.executeMethodsPerFlag(GetFlag1{&g})
	default:
		return  model.FlagNoExist.MarshallErrorAndAddFlag()
	}
}

func (g Get) executeMethodsPerFlag(test2 FlagMethods) []byte {

	if !test2.marshalBytes() {
		return model.ErrorMarshal.MarshallErrorAndAddFlag()
	}
	error := test2.checkParameters()
	if !error.IsNull(){
		return error.MarshallErrorAndAddFlag()
	}
	return test2.databaseInteraction()
}
func (g GetFlag1) marshalBytes() bool {
	err := json.Unmarshal(g.get.message, &g.get.request)
	if err != nil {
		return false
	}
	return true
}
func (g GetFlag1) checkParameters() model.Error {

	if len(g.get.request.StoneIDs) == 0 {
		return model.MissingStoneID
	}

	if len(g.get.request.Types) == 0 {
		return model.MissingType
	}
	g.get.request.Types = checkUnknownAndDuplicatedTypes(g.get.request.Types)

	return model.NoError
}

func (g GetFlag1) databaseInteraction() []byte {
	response, err := database.Select(g.get.request)
	if !err.IsNull() {
		return err.MarshallErrorAndAddFlag()
	}

	responseBytes, error := json.Marshal(response)
	if error != nil {
		return model.Error{1, error.Error()}.MarshallErrorAndAddFlag()
	}

	resultCode := make([]byte, 4)
	binary.LittleEndian.PutUint32(resultCode, 1)
	return append(resultCode, responseBytes...)
}
