package game

// Play provides the runtime with scripted callbacks.
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
	// FUTURE: replace class by an index?
	// ALT: broadcast an event to an optionally query/filtered set of objects.
	Query(class string, exact bool) IQuery
	// StopHere.
	// FIX: move into the event object, possible via callback parameter injection
	StopHere()
	// Go queue the passed runtime phrases for future execution.
	Go(RuntimePhrase, ...RuntimePhrase) IPromise
	// Random returns a number ranging from 0 to n, not including n.
	Random(n int) int
}

// RuntimePhrases are the workhorse of named actions.
// unlike named actions, they allow for typed parameters, and they dont raise events.
// in the future, it might be nice to associate a phrase with every action,
// they point towards user defined functions in some sort of user interface
// with slots which map to an "action class" -- a presentation for an action.
type RuntimePhrase interface {
	Execute(Play)
}

// IPromise doesnt currently provide a standard promise interface
// No actual chaining
type IPromise interface {
	Then(Callback)
}
