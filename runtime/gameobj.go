package runtime

import (
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/util/ident"
)

//
// For the sake of text/templates, and for *setting* values,
// the game object flattens all class properties and values into a single map.
// FIX? this probably *is* the implementation of the instance's data loaded into memory for runtime usage:
// not calling into the instance info, but being the instances's interface.
type GameObject struct {
	inst *M.InstanceInfo
	RuntimeValues
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
func (gobj *GameObject) String() string {
	// FIX: can gobj be id?
	return gobj.inst.Name()
}

//
// Return the name of the object.
//
func (gobj *GameObject) Name() string {
	return gobj.inst.Name()
}
