package app

import (
	"fmt"
	C "github.com/ionous/sashimi/console"
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/net/resource"
	R "github.com/ionous/sashimi/runtime"
	"github.com/ionous/sashimi/util/ident"
	"os"
)

//
// Builds commands which get sent to the player/client.
//
type CommandOutput struct {
	id               string
	C.BufferedOutput // TEMP: implements Print() and Println()
	serial           SerialOut
	events           Events
}

//
//
//
type Included map[ident.Id]*R.GameObject

func (this Included) Include(gobj *R.GameObject) {
	this[gobj.Id()] = gobj
}

//
//
//
type Events struct {
	array []resource.Dict
}

func (this *Events) Add(name string, value interface{}) {
	this.array = append(this.array, resource.Dict{name: value})
}

func (this *Events) Flush() []resource.Dict {
	ret := this.array
	this.array = []resource.Dict{}
	return ret
}

//
//
//
type SerialOut struct {
	*ObjectSerializer
	includes Included
}

func (this *SerialOut) TryObjectRef(gobj *R.GameObject) (ret *resource.Object, okay bool) {
	if this.IsKnown(gobj) {
		ret = this.NewObjectRef(gobj)
		okay = true
	}
	return
}

func (this *SerialOut) NewObjectRef(gobj *R.GameObject) *resource.Object {
	this.includes.Include(gobj)
	return this.NewObject(resource.ObjectList{}, gobj)
}

func (this *SerialOut) Flush() Included {
	ret := this.includes
	this.includes = make(Included)
	return ret
}

//
//
//
func NewCommandOutput(id string) *CommandOutput {
	this := &CommandOutput{
		id:     id,
		serial: SerialOut{NewObjectSerializer(), Included{}},
	}
	return this
}

//
// Add a command for an actor's line of dialog.
//
func (this *CommandOutput) ActorSays(who *R.GameObject, lines []string) {
	this.flushPending()
	//this.events.Add("say", this.serial.NewObjectRef(who).SetAttr("lines", lines))
	this.events.Add("say",
		resource.ObjectList{}.
			NewObject("_display_", "_sys_").
			SetAttr("lines", lines).
			SetAttr("speaker", this.serial.NewObjectRef(who)))
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
// Log the passed message locally, doesn't generate a client command.
//
func (this *CommandOutput) Log(message string) {
	os.Stderr.WriteString(message)
}

//
// Flush the commands to the passed output.
//
func (this *CommandOutput) FlushDocument(doc resource.DocumentBuilder) {
	this.flushPending()
	// techinically, it'd be some sort of 201 location of the new frame url.
	this.FlushFrame(doc, doc.NewIncludes())

}

//
// NOTE: Both header and included may be the same list -- as is true of the first frame.
//
func (this *CommandOutput) FlushFrame(header, included resource.IBuildObjects) {
	// create a new frame
	//include all events for this new frame
	header.NewObject(this.id, "game").SetAttr("events", this.events.Flush())
	// includes the object once, after all of properties have changed.
	for _, gobj := range this.serial.Flush() {
		this.serial.SerializeObject(included, gobj, false)
	}
}

func (this *CommandOutput) changedLocation(action *M.ActionInfo, gobjs []*R.GameObject) {
	this.flushPending()
	who, where := this.serial.NewObjectRef(gobjs[1]), this.serial.NewObjectRef(gobjs[2])
	this.events.Add("set-initial-position", who.SetMeta("location", where))

}

//
// via callback via PropertyWatcher, triggered by runtime's Notify()
//
func (this *CommandOutput) propertyChanged(game *R.Game, gobj *R.GameObject, prop M.IProperty, prev, next interface{}) {
	//
	// property changes dont cause an object to be serialized
	// some other event or request is required
	//
	switch prop := prop.(type) {
	case *M.NumProperty:
		if obj, ok := this.serial.TryObjectRef(gobj); ok {
			this.events.Add("x-num", obj.SetAttr(jsonId(prop.Id()), next))
		}

	case *M.TextProperty:
		if obj, ok := this.serial.TryObjectRef(gobj); ok {
			this.events.Add("x-txt", obj.SetAttr(jsonId(prop.Id()), next))
		}

	case *M.EnumProperty:
		if obj, ok := this.serial.TryObjectRef(gobj); ok {
			prev := jsonId(prev.(ident.Id))
			next := jsonId(next.(ident.Id))
			this.events.Add("x-set", obj.SetMeta("change-states", []string{prev, next}))
		}

	case *M.RelativeProperty:
		// get the relation
		relation := game.Model.Relations[prop.Relation()]

		// get the reverse property
		other, foundOther := relation.Other(prop.Class(), prop.Id())
		if !foundOther {
			panic(fmt.Sprint("couldnt match", prop, relation))
		}

		// fire for the original object
		if obj, ok := this.serial.TryObjectRef(gobj); ok {
			this.events.Add("x-rel", obj.SetMeta("rel", jsonId(prop.Id())))
		}

		// fire for the prev object's relationships
		if gprev, ok := game.Objects[prev.(ident.Id)]; ok {
			if obj, ok := this.serial.TryObjectRef(gprev); ok {
				this.events.Add("x-rel", obj.SetMeta("rel", jsonId(other.Property)))
			}
		}

		// fire for the next object's relationships
		if gnext, ok := game.Objects[next.(ident.Id)]; ok {
			if obj, ok := this.serial.TryObjectRef(gnext); ok {
				this.events.Add("x-rel", obj.SetMeta("rel", jsonId(other.Property)))
			}
		}
	}
}

//
// Write buffered lines into the fake $lines object
//
func (this *CommandOutput) flushPending() {
	if lines := this.BufferedOutput.Flush(); len(lines) > 0 {
		// a queriable resource so that it's reocoverable, pagination?
		this.events.Add("say", resource.ObjectList{}.NewObject("_display_", "_sys_").SetAttr("lines", lines))
	}
}
