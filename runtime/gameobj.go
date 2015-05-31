package runtime

import (
	E "github.com/ionous/sashimi/event"
	M "github.com/ionous/sashimi/model"
)

//
// For the sake of text/templates, and for *setting* values,
// the game object flattens all class properties and values into a single map.
// Could potentially use this for diffing to make save files or report changes from the initial instance values.
// note: for templates, this stores choices as the value of the choice ( rather than as their property names )
//
type GameObject struct {
	info       *M.InstanceInfo
	values     RuntimeValues
	temps      TemplatePool // this isn' terrible here, but the templates could just go into the runtime values...
	dispatcher E.Dispatcher
}

//
// Map of all game objects, keyed by model instance id.
//
type GameObjects map[M.StringId]*GameObject

//
// Return the name of the object.
//
func (this *GameObject) Id() M.StringId {
	return this.info.Id()
}

//
// Return the name of the object.
//
func (this *GameObject) String() string {
	// FIX: can this be id?
	return this.info.Name()
}

//
// E.Dispatcher
//
func (this *GameObject) Listen(evt string, handler E.IListen, capture bool) {
	this.dispatcher.Listen(evt, handler, capture)
}

//
// E.Dispatcher
//
func (this *GameObject) Silence(evt string, handler E.IListen, capture bool) {
	this.dispatcher.Silence(evt, handler, capture)
}
