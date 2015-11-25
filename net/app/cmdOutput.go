package app

import (
	C "github.com/ionous/sashimi/console"
	E "github.com/ionous/sashimi/event"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/net/resource"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/util/ident"
	"os"
)

// CommandOutput records game state changes, gets sent to the player/client.
type CommandOutput struct {
	id               string
	C.BufferedOutput // TEMP: implements Print() and Println()
	serial           SerialOut
	events           *EventStream
}

// Includes all objects referenced by the CommandOutput.
type Includes map[ident.Id]meta.Instance

func (inc Includes) Include(gobj meta.Instance) {
	inc[gobj.GetId()] = gobj
}

// SerialOut is used by CommandOutput to track instances
type SerialOut struct {
	*ObjectSerializer
	includes Includes
}

// TryObjectRef only creates an object ref if the object is already known.
func (serial *SerialOut) TryObjectRef(gobj meta.Instance) (ret *resource.Object, okay bool) {
	if serial.IsKnown(gobj) {
		ret = serial.NewObjectRef(gobj)
		okay = true
	}
	return
}

func (serial *SerialOut) NewObjectRef(gobj meta.Instance) *resource.Object {
	serial.includes.Include(gobj)
	return serial.NewObject(resource.ObjectList{}, gobj)
}

func (serial *SerialOut) Flush() Includes {
	ret := serial.includes
	serial.includes = make(Includes)
	return ret
}

// NewCommandOutput
func NewCommandOutput(id string) *CommandOutput {
	out := &CommandOutput{
		id:     id,
		serial: SerialOut{NewObjectSerializer(), Includes{}},
		events: NewEventStream(),
	}
	return out
}

// ActorSays adds a command for an actor's line of dialog.
func (out *CommandOutput) ActorSays(who meta.Instance, lines []string) {
	out.flushPending()
	tgt := out.serial.NewObjectRef(who)
	out.events.AddAction("say", tgt, lines)
}

// ScriptSays adds a command for passed script lines.
// ( The implementation actually consolidates consecutive script says into a single command.
// which gets written during flushPending() )
func (out *CommandOutput) ScriptSays(lines []string) {
	for _, l := range lines {
		out.Println(l)
	}
}

func (out *CommandOutput) BeginEvent(tgt E.ITarget, _ E.PathList, msg *E.Message) api.IEndEvent {
	out.flushPending()
	// msg.Data == RunTimeAction
	// theres not really parameters for events right now
	// other than tgt, src, ctx right now.
	out.events.PushEvent(msg.Id, tgt, nil)
	return out
}

func (out *CommandOutput) EndEvent() {
	out.flushPending()
	out.events.PopEvent()
}

// Log the passed message locally, doesn't generate a client command.
// FIX: a log interface -- perhaps as an anonymous member -- then we could have logf, etc.
func (out *CommandOutput) Log(message string) {
	os.Stderr.WriteString(message)
}

// FlushDocument containing all commands to the passed document builder.
func (out *CommandOutput) FlushDocument(doc resource.DocumentBuilder) {
	out.flushPending()
	out.flushFrame(doc, doc.NewIncludes())
}

// FlushFrame NOTE: Both header and included may be the same list -- as is true of the first frame.
func (out *CommandOutput) flushFrame(header, included resource.IBuildObjects) {
	// create a new frame
	// include all events for out new frame
	game := header.NewObject(out.id, "game")
	if events := out.events.Flush(); len(events) > 0 {
		game.SetAttr("events", events)
	}
	// includes the object once, after all of properties have changed.
	for _, gobj := range out.serial.Flush() {
		out.serial.SerializeObject(included, gobj, false)
	}
}

func (out *CommandOutput) NumChange(gobj meta.Instance, prop ident.Id, prev, next float32) {
	if obj, ok := out.serial.TryObjectRef(gobj); ok {
		data := struct {
			Prop  string  `json:"prop"`
			Value float32 `json:"value"`
		}{jsonId(prop), next}
		out.events.AddAction("x-num", obj, data)
	}
}

func (out *CommandOutput) TextChange(gobj meta.Instance, prop ident.Id, prev, next string) {
	if obj, ok := out.serial.TryObjectRef(gobj); ok {
		data := struct {
			Prop  string `json:"prop"`
			Value string `json:"value"`
		}{jsonId(prop), next}
		out.events.AddAction("x-txt", obj, data)
	}
}
func (out *CommandOutput) StateChange(gobj meta.Instance, prop ident.Id, prev, next ident.Id) {
	if obj, ok := out.serial.TryObjectRef(gobj); ok {
		data := struct {
			Prop string `json:"prop"`
			Prev string `json:"prev"`
			Next string `json:"next"`
		}{jsonId(prop),
			jsonId(prev),
			jsonId(next)}
		out.events.AddAction("x-set", obj, data)
	}
}
func (out *CommandOutput) ReferenceChange(gobj meta.Instance, prop, other ident.Id, prev, next meta.Instance) {
	if out.serial.IsKnown(gobj) || out.serial.IsKnown(prev) || out.serial.IsKnown(next) {
		obj := out.serial.NewObjectRef(gobj)
		relChange := struct {
			Prop  string           `json:"prop"`
			Other string           `json:"other"`
			Prev  *resource.Object `json:"prev,omitempty"`
			Next  *resource.Object `json:"next,omitempty"`
		}{Prop: jsonId(prop), Other: jsonId(other)}

		// fire for the prev object's relationships
		if prev != nil {
			relChange.Prev = out.serial.NewObjectRef(prev)
		}

		// fire for the next object's relationships
		if next != nil {
			relChange.Next = out.serial.NewObjectRef(next)
		}
		out.events.AddAction("x-rel", obj, relChange)
	}
}

// flushPending buffered lines into the fake display object.
// ( so long as theres a flush before push and pop. )
func (out *CommandOutput) flushPending() {
	if lines := out.BufferedOutput.Flush(); len(lines) > 0 {
		// FIXFIXIX: theres some sort of bug in the buffered output or the code that uses it,
		// leading to empty, and unconsolidated, "say" staements
		// this can be seen in command_test: after lines": ["", "lab", "an empty room", ""],
		// are a series of blank says.
		empty := true
		for _, l := range lines {
			if len(l) > 0 {
				empty = false
				break
			}
		}
		if !empty {
			// FIX? a queriable resource so that it's recoverable, pagination?
			var tgt = resource.ObjectList{}.NewObject("_display_", "_sys_")
			out.events.AddAction("print", tgt, lines)
		}
	}
}
