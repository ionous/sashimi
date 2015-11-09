package model

import (
	"github.com/ionous/sashimi/model/table"
)

// Model: Compiled results of a sashimi story.
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
	EventListeners EventCallbacks
	Tables         table.Tables
	//
	NounNames      NounNames
	SingleToPlural SingleToPlural

	// this cant serialize: need to replace with
	// lists and table variables in globals
	//Generators     GeneratorMap
}

type SingleToPlural map[string]string
