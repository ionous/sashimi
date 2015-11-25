package api

import (
	"github.com/ionous/sashimi/meta"
)

// Output required to display text.
type Output interface {
	// ScriptSays some standard output.
	ScriptSays(lines []string)
	// ActorSays that the actor or object with the passed name has something to say.
	ActorSays(whose meta.Instance, lines []string)
	// Log debugging output.
	Log(string)
	// Println is used by StandardGame.Input to display the results of bad parsing to the user.
	// FIX: merge into some "report" interface for scripts; logging, printing, reporting need to be thoroughly cleaned up. categories for the logging, possibly self registerable, enumerable, listable would be great.
	Println(...interface{})
}
