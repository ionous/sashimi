package runtime

//
// interface the runtime requires an implementation of in order to display text.
//
type IOutput interface {
	//
	// Standard output.
	//
	ScriptSays(lines []string)
	//
	// The actor or object with the passed name has something to say.
	//
	ActorSays(whose *GameObject, lines []string)
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
