package simple

import (
	"fmt"
	C "github.com/ionous/sashimi/console"
	"github.com/ionous/sashimi/runtime/api"
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
func (this *SimpleOutput) ActorSays(whose api.Instance, lines []string) {
	fmt.Println("Actor says", lines)
	prop, _ := whose.GetProperty("Name")
	name := prop.GetValue().GetText()
	for _, l := range lines {
		this.Println(name, ": ", l)
	}
}

//
func (this *SimpleOutput) Log(out string) {
	os.Stderr.WriteString(out)
}
