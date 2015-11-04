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
	NounNames      NounNames
	ActionHandlers ActionCallbacks
	EventListeners ListenerCallbacks
	Tables         table.Tables
	//Generators     GeneratorMap
}
