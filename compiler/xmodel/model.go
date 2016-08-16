package xmodel

import (
	"github.com/ionous/sashimi/compiler/model/table"
)

// Model the compiler's original output
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
}

type SingleToPlural map[string]string
