package game

type SourceLookup func(Play) IObject
type TargetLookup func(Play, IObject) IObject
type ContextLookup func(Play, IObject, IObject) IObject

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
	// NewFrom: creates a new object from the passed plural-named class.
	NewFrom(data string) IObject
	// Say friendly narrative print.
	Say(text ...string)
	// Log a quiet print.
	Log(text ...interface{})
	// Visit finds one or more objects by class name.
	// Return true from the passed function to terminate the search.
	// FIX: replace by a query?
	// ALT: broadcast an event to an optionally query/filtered set of objects.
	Visit(class string, visits func(IObject) bool) bool
	//
	Rules() IGameRules
	// StopHere.
	// FIX: move into the event object, possible via callback parameter injection
	StopHere()
	// FIX: a hack, mainly for conversations
	// a system of varient -- and possibly user type -- globals is needed
	// but many of these could go away if there was a real table implementatioin
	Global(name string) (interface{}, bool)
	Go(RuntimePhrase)
	// return a random number ranging from 0 to n, not including n
	Random(n int) int
}

type RuntimePhrase interface {
	Execute(Play)
}

// Any returns the first compatible class found.
func Any(g Play, class string) (ret IObject, found bool) {
	found = g.Visit(class, func(obj IObject) bool {
		ret = obj
		return true
	})
	return ret, found
}
