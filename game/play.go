package game

type SourceLookup func(Play) IObject
type TargetLookup func(Play, IObject) IObject
type ContextLookup func(Play, IObject, IObject) IObject
type Callback func(Play)

type IGameRules interface {
	PushParserSource(SourceLookup)
	PopParserSource()

	PushParentLookup(TargetLookup)
	PopParentLookup()
}

// FUTURE: replace this interface with a set of global functions which delegate based on context:
// the script system for definitions, the game systems for callbacks.
// augment the callbacks with dependency injection to provide standard objects.
type Play interface {
	The(noun string) IObject
	Our(noun string) IObject
	A(noun string) IObject
	// friendly narrative print
	Say(text ...string)
	// system version of say
	//Report(text ...string)
	// quieter version of report
	Log(text ...string)
	// find some object by class name
	// it's interesting how similar that is to finding an instance via a relation....
	Any(class string) IObject
	//
	Rules() IGameRules
	// FIX: move into the event object, possible via callback parameter injection
	StopHere()
}
