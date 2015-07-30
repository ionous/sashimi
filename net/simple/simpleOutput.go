package simple

import (
	"fmt"
	C "github.com/ionous/sashimi/console"
	R "github.com/ionous/sashimi/runtime"
	"os"
)

// implements IOutput
type SimpleOutput struct {
	C.BufferedOutput // implements Print() and Println()
}

//
func (this *SimpleOutput) ScriptSays(lines []string) {
	fmt.Println("Script says", lines)
	for _, l := range lines {
		this.Println(l)
	}
}

//
func (this *SimpleOutput) ActorSays(whose *R.GameObject, lines []string) {
	fmt.Println("Actor says", lines)
	name := whose.Value("Name")
	for _, l := range lines {
		this.Println(name, ": ", l)
	}
}

//
func (this *SimpleOutput) Log(out string) {
	os.Stderr.WriteString(out)
}
