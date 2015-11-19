package simple

import (
	"fmt"
	C "github.com/ionous/sashimi/console"
	"github.com/ionous/sashimi/meta"
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
func (this *SimpleOutput) ActorSays(whose meta.Instance, lines []string) {
	fmt.Println("Actor says", lines)
	var name string
	if prop, ok := whose.FindProperty("name"); ok {
		name = prop.GetValue().GetText()
	}
	for _, l := range lines {
		this.Println(name, ": ", l)
	}
}

//
func (this *SimpleOutput) Log(out string) {
	os.Stderr.WriteString(out)
}
