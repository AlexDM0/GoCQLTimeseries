package parser

import (
	"CrownstoneServer/database"
	"CrownstoneServer/model"
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
	message    []byte
	request		*model.RequestSelectJSON
}




func (g Get) parseFlag() bool {
	flag := model.ByteToInt(g.message, 0)
	switch flag {
	case 1:
		return true
	default:
		return false
	}
}
func (g Get) parseJSON() bool {
	err := json.Unmarshal(g.message[4:], &g.request)
	if err != nil {
		return false
	}
	return true
}
func (g Get) checkParameters() model.Error {

	if len(g.request.StoneIDs) == 0{
		return model.Error{100, "StoneID is missing"}
	}

	if len(g.request.Types) == 0 {
		return model.Error{100, "Type is missing"}
	}
	g.request.Types = checkUnknownAndDuplicatedTypes(g.request.Types)

	return model.Null
}

func (g Get) parseJSONToDatabaseQueries() []byte {
	response, err := database.Select(g.request)
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

