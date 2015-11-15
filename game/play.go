package game

type SourceLookup func(Play) IObject
type TargetLookup func(Play, IObject) IObject
type ContextLookup func(Play, IObject, IObject) IObject

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
	// Say friendly narrative print.
	Say(text ...string)
	// Log a quiet print.
	Log(text ...interface{})
	// List one or more objects by class name.
	// FIX: replace by a query?
	// ALT: broadcast an event to an optionally query/filtered set of objects.
	List(class string) IList
	// StopHere.
	// FIX: move into the event object, possible via callback parameter injection
	StopHere()
	//
	Go(RuntimePhrase)
	// return a random number ranging from 0 to n, not including n
	Random(n int) int
}

type RuntimePhrase interface {
	Execute(Play)
}
