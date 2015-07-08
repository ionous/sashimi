package runtime

import (
	E "github.com/ionous/sashimi/event"
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/util/ident"
)

//
// For the sake of text/templates, and for *setting* values,
// the game object flattens all class properties and values into a single map.
// Could potentially use gobj for diffing to make save files or report changes from the initial instance values.
// note: for templates, gobj stores choices as the value of the choice ( rather than as their property names )
//
type GameObject struct {
	inst *M.InstanceInfo
	RuntimeValues
	temps      TemplatePool // gobj isn' terrible here, but the templates could just go into the runtime values...
	dispatcher E.Dispatcher
}

//
// Map of all game objects, keyed by model instance id.
//
type GameObjects map[ident.Id]*GameObject

//
// Return the name of the object.
//
func (gobj *GameObject) Id() ident.Id {
	return gobj.inst.Id()
}

//
//
//
func (gobj *GameObject) Class() *M.ClassInfo {
	return gobj.inst.Class()
}

//
// Return the name of the object.
//
func (gobj *GameObject) Name() string {
	return gobj.inst.Name()
}

//
// Return the name of the object.
//
func (gobj *GameObject) String() string {
	// FIX: can gobj be id?
	return gobj.inst.Name()
}

//
// E.Dispatcher
//
func (gobj *GameObject) Listen(evt string, handler E.IListen, capture bool) {
	gobj.dispatcher.Listen(evt, handler, capture)
}

//
// E.Dispatcher
//
func (gobj *GameObject) Silence(evt string, handler E.IListen, capture bool) {
	gobj.dispatcher.Silence(evt, handler, capture)
}
