package api

// Output requires an implementation to display text.
type Output interface {
	//
	// Standard output.
	//
	ScriptSays(lines []string)
	//
	// The actor or object with the passed name has something to say.
	//
	ActorSays(whose Instance, lines []string)
	//
	// Debugging output.
	//
	Log(string)
	//
	// FIXFIXFIXFIX
	// this is used by StandardGame.Input to display the results of bad parsing to the user
	// merge into some "report"
	//
	Println(...interface{})
}
