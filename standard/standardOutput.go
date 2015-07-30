package standard

import (
	"fmt"
	C "github.com/ionous/sashimi/console"
	R "github.com/ionous/sashimi/runtime"
	"io"
	"strings"
)

type StandardOutput struct {
	console   C.IConsole
	writer    io.Writer
	lastEmpty bool // collapse mutliple empty lines
}

func (out *StandardOutput) Println(args ...interface{}) {
	str := fmt.Sprint(args...)
	nowEmpty := len(strings.TrimSpace(str)) == 0
	if !nowEmpty {
		if out.lastEmpty {
			out.console.Println(" ")
		}
		out.console.Println(str)
	}
	out.lastEmpty = nowEmpty
}

func (out *StandardOutput) ScriptSays(lines []string) {
	for _, l := range lines {
		out.Println(l)
	}
}

func (out *StandardOutput) ActorSays(whose *R.GameObject, lines []string) {
	if len(lines) > 0 {
		// in other contexts ActorSays needs R.GameObject for SerializeObject
		name := whose.Value("Name").(string)
		out.console.Println(name, ": ", lines[0])
		for _, l := range lines[1:] {
			out.Println(l)
		}
	}
}

func (out *StandardOutput) Log(s string) {
	out.writer.Write([]byte(s))
}
