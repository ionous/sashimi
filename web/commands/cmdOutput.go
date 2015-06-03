package commands

import (
	C "github.com/ionous/sashimi/console"
	"os"
)

//
// Builds commands which get sent to the player/client.
//
type CommandOutput struct {
	id               string
	cmds             []Dict // every command is an json-like object
	C.BufferedOutput        // TEMP: implements Print() and Println()
}

//
// Add the named command and its pased json-like data.
//
func (this *CommandOutput) NewCommand(name string, data Dict) {
	this.flushPending()
	cmd := Dict{name: data}
	this.cmds = append(this.cmds, cmd)
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
	this.NewCommand("say", Dict{"actor": name, "lines": lines})
}

//
// Log the passed message locally, doesn't generate a client command.
//
func (this *CommandOutput) Log(message string) {
	os.Stderr.WriteString(message)
}

//
// Flush the commands to the passed output.
//
func (this *CommandOutput) Fetch() (data interface{}, err error) {
	this.flushPending()
	cmds := Dict{"id": this.id, "cmds": this.cmds}
	this.cmds = nil
	return cmds, err
}

//
// Write buffered lines as a single command "say".
//
func (this *CommandOutput) flushPending() {
	if lines := this.Flush(); len(lines) > 0 {
		this.NewCommand("say", Dict{"lines": lines})
	}
}
