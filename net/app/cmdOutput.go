package app

import (
	C "github.com/ionous/sashimi/console"
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/net/resource"
	R "github.com/ionous/sashimi/runtime"
	"github.com/ionous/sashimi/util/ident"
	"os"
)

//
// CommandOutput records game state changes, gets sent to the player/client.
//
type CommandOutput struct {
	id               string
	C.BufferedOutput // TEMP: implements Print() and Println()
	serial           SerialOut
	events           *EventStream
}

//
// Includes all objects referenced by the CommandOutput.
//
type Includes map[ident.Id]*R.GameObject

func (inc Includes) Include(gobj *R.GameObject) {
	inc[gobj.Id()] = gobj
}

//
// SerialOut is used by CommandOuput to track runtime.GameObjects.
//
type SerialOut struct {
	*ObjectSerializer
	includes Includes
}

// TryObjectRef only creates an object ref if the object is already known.
func (serial *SerialOut) TryObjectRef(gobj *R.GameObject) (ret *resource.Object, okay bool) {
	if serial.IsKnown(gobj.Id()) {
		ret = serial.NewObjectRef(gobj)
		okay = true
	}
	return
}

func (serial *SerialOut) NewObjectRef(gobj *R.GameObject) *resource.Object {
	serial.includes.Include(gobj)
	return serial.NewObject(resource.ObjectList{}, gobj)
}

func (serial *SerialOut) Flush() Includes {
	ret := serial.includes
	serial.includes = make(Includes)
	return ret
}

//
// NewCommandOutput
//
func NewCommandOutput(id string) *CommandOutput {
	out := &CommandOutput{
		id:     id,
		serial: SerialOut{NewObjectSerializer(), Includes{}},
		events: NewEventStream(),
	}
	return out
}

//
// ActorSays adds a command for an actor's line of dialog.
//
func (out *CommandOutput) ActorSays(who *R.GameObject, lines []string) {
	out.flushPending()
	tgt := out.serial.NewObjectRef(who)
	out.events.AddAction("say", tgt, lines)
}

//
// ScriptSays add a command for passed script lines.
// ( The implementation actually consolidates consecutive script says into a single command. )
//
func (out *CommandOutput) ScriptSays(lines []string) {
	for _, l := range lines {
		out.Println(l)
	}
}

//
// Log the passed message locally, doesn't generate a client command.
// FIX: a log interface -- perhaps as an anonymous member -- then we could have logf, etc.
func (out *CommandOutput) Log(message string) {
	os.Stderr.WriteString(message)
}

//
// FlushDocument containing all commands to the passed document builder.
//
func (out *CommandOutput) FlushDocument(doc resource.DocumentBuilder) {
	out.flushPending()
	// techinically, it'd be some sort of 201 location of the new frame url.
	out.flushFrame(doc, doc.NewIncludes())
}

//
// FlushFrame NOTE: Both header and included may be the same list -- as is true of the first frame.
//
func (out *CommandOutput) flushFrame(header, included resource.IBuildObjects) {
	// create a new frame
	//include all events for out new frame
	game := header.NewObject(out.id, "game")
	if events := out.events.Flush(); len(events) > 0 {
		game.SetAttr("events", events)
	}
	// includes the object once, after all of properties have changed.
	for _, gobj := range out.serial.Flush() {
		out.serial.SerializeObject(included, gobj, false)
	}
}

func (this *CommandOutput) changedLocation(action *M.ActionInfo, gobjs []*R.GameObject) {
	this.flushPending()
	who, where := this.serial.NewObjectRef(gobjs[1]), this.serial.NewObjectRef(gobjs[2])
	this.events.AddAction("set-initial-position", who, where)
}

// propertyChanged callback via PropertyWatcher, triggered by runtime's Notify()
func (out *CommandOutput) propertyChanged(game *R.Game, gobj *R.GameObject, prop M.IProperty, prev, next interface{}) {
	//
	// property changes dont cause an object to be serialized
	// some other event or request is required
	//
	switch prop := prop.(type) {
	case M.NumProperty:
		if obj, ok := out.serial.TryObjectRef(gobj); ok {
			data := struct {
				Prop  string  `json:"prop"`
				Value float32 `json:"value"`
			}{jsonId(prop.GetId()), next.(float32)}
			out.events.AddAction("x-num", obj, data)
		}

	case M.TextProperty:
		if obj, ok := out.serial.TryObjectRef(gobj); ok {
			data := struct {
				Prop  string `json:"prop"`
				Value string `json:"value"`
			}{jsonId(prop.GetId()), next.(string)}
			out.events.AddAction("x-txt", obj, data)
		}

	case M.EnumProperty:
		if obj, ok := out.serial.TryObjectRef(gobj); ok {
			data := struct {
				Prop string `json:"prop"`
				Prev string `json:"prev"`
				Next string `json:"next"`
			}{jsonId(prop.GetId()),
				jsonId(prev.(ident.Id)),
				jsonId(next.(ident.Id))}
			out.events.AddAction("x-set", obj, data)
		}

	case M.RelativeProperty:
		// get the relation
		relation := game.Model.Relations[prop.Relation]

		// get the reverse property
		other := relation.GetOther(prop.IsRev)

		type RelationChange struct {
			Prop  string           `json:"prop"`
			Other string           `json:"other"`
			Prev  *resource.Object `json:"prev,omitempty"`
			Next  *resource.Object `json:"next,omitempty"`
		}

		// fire for the original object
		if out.serial.IsKnown(gobj.Id()) || out.serial.IsKnown(prev.(ident.Id)) || out.serial.IsKnown(next.(ident.Id)) {
			obj := out.serial.NewObjectRef(gobj)
			relChange := RelationChange{Prop: jsonId(prop.GetId()), Other: jsonId(other.Property)}

			// fire for the prev object's relationships
			if gprev, ok := game.Objects[prev.(ident.Id)]; ok {
				relChange.Prev = out.serial.NewObjectRef(gprev)
			}

			// fire for the next object's relationships
			if gnext, ok := game.Objects[next.(ident.Id)]; ok {
				relChange.Next = out.serial.NewObjectRef(gnext)
			}
			out.events.AddAction("x-rel", obj, relChange)
		}
	}
}

// flushPending buffered lines into the fake display object.
func (out *CommandOutput) flushPending() {
	// only if theres a flush before push and pop.
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
