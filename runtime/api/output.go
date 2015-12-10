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
}
