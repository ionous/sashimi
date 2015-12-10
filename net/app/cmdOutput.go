package app

import (
	"fmt"
	"github.com/ionous/sashimi/console"
	E "github.com/ionous/sashimi/event"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/net/resource"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/util/ident"
	"os"
)

// CommandOutput records game state changes, gets sent to the player/client.
type CommandOutput struct {
	id     string
	view   View
	events *EventStream
	serial *ObjSerializer
	text   console.BufferedOutput
}

func NewCommandOutput(id string, m meta.Model, view View) *CommandOutput {
	return &CommandOutput{
		id:     id,
		view:   view,
		events: NewEventStream(),
		serial: NewObjSerializer(m, resource.NewObjectList()),
	}
}

// ActorSays adds a command for an actor's line of dialog.
func (out *CommandOutput) ActorSays(gobj meta.Instance, lines []string) {
	if !out.view.InView(gobj) {
		out.Log(fmt.Sprintf("CommandOutput: actor '%s' not in view, ignoring speech: (%v)\n", gobj.GetId(), lines))
	} else {
		out.flushPending()
		tgt := NewObjectRef(gobj)
		out.events.AddAction("say", tgt, lines)
	}
}

// ScriptSays adds a command for passed script lines.
// ( The implementation actually consolidates consecutive script says into a single command.
// which gets written during flushPending() )
func (out *CommandOutput) ScriptSays(lines []string) {
	for _, l := range lines {
		out.text.Println(l)
	}
}

func (out *CommandOutput) FlushFrame() {
	out.flushPending()
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
	game := doc.NewObject(out.id, "game")
	if events := out.events.Flush(); len(events) > 0 {
		game.SetAttr("events", events)
	}
	doc.SetIncluded(out.serial.out)
}

func (out *CommandOutput) NumChange(gobj meta.Instance, prop ident.Id, prev, next float32) {
	if !out.view.InView(gobj) {
		out.Log(fmt.Sprintf("CommandOutput: '%s' not in view,ignoring num change %s(%s->%s)\n", gobj.GetId(), prop, prev, next))
	} else {
		obj := NewObjectRef(gobj)
		data := struct {
			Prop  string  `json:"prop"`
			Value float32 `json:"value"`
		}{jsonId(prop), next}
		out.events.AddAction("x-num", obj, data)
	}
}

func (out *CommandOutput) TextChange(gobj meta.Instance, prop ident.Id, prev, next string) {
	if !out.view.InView(gobj) {
		out.Log(fmt.Sprintf("CommandOutput: '%s' not in view(%s), ignoring text change %s(%s->%s)\n", gobj.GetId(), prop, prev, next))
	} else {
		obj := NewObjectRef(gobj)
		data := struct {
			Prop  string `json:"prop"`
			Value string `json:"value"`
		}{jsonId(prop), next}
		out.events.AddAction("x-txt", obj, data)
	}
}

func (out *CommandOutput) StateChange(gobj meta.Instance, prop ident.Id, prev, next ident.Id) {
	if !out.view.InView(gobj) {
		out.Log(fmt.Sprintf("CommandOutput: '%s' not in view, ignoring state change %s(%s->%s)\n", gobj.GetId(), prop, prev, next))
	} else {
		obj := NewObjectRef(gobj)
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

// currently sent on the "one" side of an object for any object value.
// ( ex. actor.whereabouts; not: room.contents. )
func (out *CommandOutput) ReferenceChange(gobj meta.Instance, prop, other ident.Id, prev, next meta.Instance) {
	if out.view.Viewpoint() == gobj {
		var n ident.Id
		if next != nil {
			n = next.GetId()
		}
		out.Log(fmt.Sprintf("CommandOutput: changing view to %s\n", n))
		if out.view.ChangedView(gobj, prop, next) && next != nil {
			out.serial.Include(next)
		}
	} else {
		if out.view.EnteredView(gobj, prop, next) {
			out.serial.Include(gobj)
		}
	}
	relatedView := out.view.InView(gobj) ||
		(prev != nil && out.view.InView(prev)) ||
		(next != nil && out.view.InView(next))

	if !relatedView {
		var p, n ident.Id
		if prev != nil {
			p = prev.GetId()
		}
		if next != nil {
			n = next.GetId()
		}
		out.Log(fmt.Sprintf("CommandOutput: '%s' not in view, ignoring refchange %v(%v->%v)\n", gobj.GetId(), prop, p, n))
	} else {
		obj := NewObjectRef(gobj)

		relChange := struct {
			Prop  string           `json:"prop"`
			Other string           `json:"other"`
			Prev  *resource.Object `json:"prev,omitempty"`
			Next  *resource.Object `json:"next,omitempty"`
		}{Prop: jsonId(prop), Other: jsonId(other)}

		if prev != nil {
			relChange.Prev = NewObjectRef(prev)
		}
		if next != nil {
			relChange.Next = NewObjectRef(next)
		}

		out.events.AddAction("x-rel", obj, relChange)
	}
}

// flushPending buffered lines into the fake display object.
// ( so long as theres a flush before push and pop. )
func (out *CommandOutput) flushPending() {
	if lines := out.text.Flush(); len(lines) > 0 {
		if !EmptyLines(lines) {
			tgt := resource.NewObject("_display_", "_sys_")
			out.events.AddAction("print", tgt, lines)
		}
	}
}

// FXIXI some code -- like report the view --
// prettifies the output by printing blank lines (ugh)
func EmptyLines(lines []string) bool {
	empty := true
	for _, l := range lines {
		if len(l) > 0 {
			empty = false
			break
		}
	}
	return empty
}
