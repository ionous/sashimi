package simple

import (
	"fmt"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util"
	"os"
)

// implements IOutput
type SimpleOutput struct {
	util.BufferedOutput // implements Print() and Println()
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
