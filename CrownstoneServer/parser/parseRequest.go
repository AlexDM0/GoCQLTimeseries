package parser

import (
	"CrownstoneServer/model"
	"strings"
)

type Command interface {
	parseFlag() bool
	parseJSON() bool
	checkParameters() model.Error
	parseJSONToDatabaseQueries() []byte
}

func ParseOpCode(opCode uint32, message []byte) []byte {
	switch opCode {
	//TODO: insert
	case 100:
		i := Insert{message, &model.InsertJSON{}}
		return parser(i)
	case 200:
		s := Get{message, &model.RequestSelectJSON{}}
		return parser(s)
	case 500:
		d := Delete{message, &model.DeleteJSON{}}
		return parser(d)
		//TODO: Research delete management
	default:
		return model.Error{10, "Server doesn't recognise opcode"}.MarshallErrorAndAddFlag()
	}
}

func parser(c Command) []byte {
	err := c.parseFlag()
	if !err {
		return model.Error{0, "Flag doesn't exist"}.MarshallErrorAndAddFlag()
	}
	err = c.parseJSON()
	if !err {
		return model.Error{0, "Problem with parsing JSON"}.MarshallErrorAndAddFlag()
	}
	errBytes := c.checkParameters()
	if !errBytes.IsNull() {
		return errBytes.MarshallErrorAndAddFlag()
	}

	result := c.parseJSONToDatabaseQueries()
	return result
}

func checkUnknownAndDuplicatedTypes(request []string) ([]string) {
	var typeList = []bool{false, false, false}
	for _, v := range request {
		switch strings.ToLower(v) {
		case model.UnitW:
			typeList[0] = true
		case model.Unitpf:
			typeList[1] = true
		case model.UnitkWh:
			typeList[2] = true
		}
	}
	typePerQuery := make([]string, 2)
	if typeList[0] && typeList[1] {
		typePerQuery[0] = model.UnitWAndpf
	} else if typeList[0] {
		typePerQuery[0] = model.UnitW
	} else if typeList[1] {
		typePerQuery[0] = model.Unitpf
	}

	if typeList[2] {
		typePerQuery[1] = model.UnitkWh
	}

	return typePerQuery

}