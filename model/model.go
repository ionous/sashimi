package model

import (
	"encoding/json"
	"github.com/ionous/sashimi/model/table"
	"github.com/ionous/sashimi/util/ident"
)

// Model: Compiled results of a sashimi story.
type Model struct {
	Actions        Actions
	Classes        Classes
	Enumerations   Enumerations
	Events         Events
	Instances      Instances
	Aliases        Aliases
	ParserActions  []ParserAction
	Relations      Relations
	SingleToPlural SingleToPlural
	Tables         table.Tables
}

type Actions map[ident.Id]*ActionModel
type Classes map[ident.Id]*ClassModel
type Enumerations map[ident.Id]*EnumModel
type Events map[ident.Id]*EventModel
type Instances map[ident.Id]*InstanceModel
type Relations map[ident.Id]*RelationModel
type SingleToPlural map[string]string

func (m *Model) PrintModel(printer func(...interface{})) (err error) {
	if prettyBytes, e := json.MarshalIndent(m, "", " "); e != nil {
		err = e
	} else {
		printer(string(prettyBytes))
	}
	return err
}
