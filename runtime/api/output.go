package api

import (
	"github.com/ionous/sashimi/meta"
)

// Output required to display text.
// FIX: this should be removed
// ScriptSays and ActorSays should be events
// the console client should use game events to output to a standard string stream
// TBD: evaluate for best practices multiple events ( log, script, actor ) vs. one event and type info
// TBD: logging. less context passing, more java style: private final static Logger Log = LoggerFactory.getLogger(Class);
// how to accomplish that well in go? just use log?
type Output interface {
	// ScriptSays some standard output.
	ScriptSays(lines []string)
	// ActorSays that the actor or object with the passed name has something to say.
	ActorSays(whose meta.Instance, lines []string)
	// Log debugging output.
	Log(string)
}
