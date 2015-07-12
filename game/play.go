package game

type SourceLookup func(Play) IObject
type TargetLookup func(Play, IObject) IObject
type ContextLookup func(Play, IObject, IObject) IObject
type Callback func(Play)

type IGameRules interface {
	PushParentLookup(TargetLookup)
	PopParentLookup()
}

// Play provides an interface to the runtime for scripted callbacks.
// FUTURE? replace this interface with a set of global functions which delegate based on context:
// the script system for definitions, the game systems for callbacks.
// augment the callbacks with dependency injection to provide standard objects.
type Play interface {
	// The function retrieves a script declared instance.
	The(noun string) IObject
	// Our alias for The.
	Our(noun string) IObject
	// A(n) alias for The.
	A(noun string) IObject
	// New an object at runtime; can only be used with "data" classes
	Add(data string) IObject
	// Remove a previously new'd data object.
	Remove(IObject)
	// Say friendly narrative print.
	Say(text ...string)
	// Log a quiet print.
	Log(text ...string)
	// Any finds an object by class name.
	// FUTURE: replace by a more generic query which supports data as well...
	// A query will either require a cursor/iterator or a list.
	Any(class string) IObject
	//
	Rules() IGameRules
	// StopHere.
	// FIX: move into the event object, possible via callback parameter injection
	StopHere()
}
