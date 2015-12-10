package framework

import (
	"fmt"
	C "github.com/ionous/sashimi/console"
	"github.com/ionous/sashimi/meta"
	"io"
	"strings"
	"unicode"
)

type StandardOutput struct {
	console        C.IConsole
	logger         io.Writer
	lastEmpty      bool   // collapse mutliple empty lines
	lastActor      string // collapse multiple speaker lines
	multiLineActor bool
}

func NewStandardOutput(c C.IConsole, logger io.Writer) *StandardOutput {
	return &StandardOutput{console: c, logger: logger}
}
func (out *StandardOutput) Println(args ...interface{}) {
	str := fmt.Sprint(args...)
	nowEmpty := strings.TrimRightFunc(str, unicode.IsSpace) == ""
	if !nowEmpty {
		if out.lastEmpty || out.multiLineActor {
			out.console.Println(" ")
		}
		out.console.Println(str)
	}
	out.lastEmpty = nowEmpty
	out.lastActor = ""
	out.multiLineActor = false
}

func (out *StandardOutput) ScriptSays(lines []string) {
	for _, l := range lines {
		out.Println(l)
	}
}

func (out *StandardOutput) ActorSays(whose meta.Instance, lines []string) {
	if len(lines) > 0 {
		// in other contexts ActorSays needs the instance for SerializeObject
		// FIX: what about proper name?
		if prop, ok := whose.FindProperty("name"); !ok {
			panic("actor is unamed")
		} else {
			name := prop.GetValue().GetText()

			if out.lastActor != name {
				if out.lastActor != "" {
					out.lastEmpty = true
				}
				out.Println(name, ": ", lines[0])
			} else {
				out.Println(lines[0])
			}
			for _, l := range lines[1:] {
				out.Println(l)
			}
			// tricky hack: since out.println overwrites this, we set it last.
			out.lastActor = name
			out.multiLineActor = len(lines) > 1
		}
	}
}

func (out *StandardOutput) Log(s string) {
	out.logger.Write([]byte(s))
}
