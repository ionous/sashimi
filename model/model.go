package model

import (
	"github.com/ionous/sashimi/model/table"
)

//
// Results of compilation.
//
type Model struct {
	// rule like:
	Classes       ClassMap
	Relations     RelationMap
	Actions       ActionMap
	Events        EventMap
	ParserActions []ParserAction
	// data like:
	Instances      InstanceMap
	ActionHandlers ActionCallbacks
	EventListeners ListenerCallbacks
	Tables         table.Tables
	//
	NounNames      NounNames
	SingleToPlural SingleToPlural

	// this cant serialize: nned to replace with
	// lists and table variables in globals
	//Generators     GeneratorMap
}

type SingleToPlural map[string]string
