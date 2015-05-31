package commands

import (
	"encoding/json"
	C "github.com/ionous/sashimi/console"
	"io"
	"os"
)

//
// Builds commands which get sent to the player/client.
//
type CommandOutput struct {
	cmds             []Dict // every command is an json-like object
	C.BufferedOutput        // TEMP: implements Print() and Println()
}

//
// Add the named command and its pased json-like data.
//
func (this *CommandOutput) Add(name string, data Dict) {
	this.flushPending()
	cmd := Dict{name: data}
	this.cmds = append(this.cmds, cmd)
}

//
// Log the passed message locally, doesn't generate a client command.
//
func (this *CommandOutput) Log(message string) {
	os.Stderr.WriteString(message)
}

//
// Add a command for passed script lines.
// ( The implementation actually consolidates consecutive script says into a single command. )
//
func (this *CommandOutput) ScriptSays(lines []string) {
	for _, l := range lines {
		this.Println(l)
	}
}

//
// Add a command for an actor's line of dialog.
//
func (this *CommandOutput) ActorSays(name string, lines []string) {
	this.Add("speak", Dict{"actor": name, "lines": lines})
}

//
// Flush the commands to the passed output.
//
func (this *CommandOutput) Write(w io.Writer) (err error) {
	this.flushPending()
	err = json.NewEncoder(w).Encode(this.cmds)
	this.cmds = nil
	return err
}

//
// Write buffered lines as a single command "say".
//
func (this *CommandOutput) flushPending() {
	if lines := this.Flush(); len(lines) > 0 {
		this.Add("say", Dict{"lines": lines})
	}
}
