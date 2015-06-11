package standard

import (
	"fmt"
	C "github.com/ionous/sashimi/console"
	R "github.com/ionous/sashimi/runtime"
	"io"
)

type StandardOutput struct {
	console C.IConsole
	writer  io.Writer
}

func (this StandardOutput) Println(args ...interface{}) {
	str := fmt.Sprintln(args...)
	this.console.Println(str)
}

func (this StandardOutput) ScriptSays(lines []string) {
	for _, l := range lines {
		this.console.Println(l)
	}
	//this.console.Println()
}

func (this StandardOutput) ActorSays(whose *R.GameObject, lines []string) {
	if len(lines) > 0 {
		name := whose.Name()
		this.console.Println(name, ": ", lines[0])
		for _, l := range lines[1:] {
			this.console.Println(l)
		}
		this.console.Println()
	}
}

func (this StandardOutput) Log(s string) {
	this.writer.Write([]byte(s))
}
