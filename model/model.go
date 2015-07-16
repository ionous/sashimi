package model

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
	Tables         TableRelations
	Generators     GeneratorMap
}
