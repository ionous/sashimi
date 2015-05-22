package runtime

import (
	"fmt"
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
func (this *GameObject) String() string {
	return this.info.Name()
}

//
// For debugging, return the model data this object was created from.
//
func (this *GameObject) Info() *M.InstanceInfo {
	return this.info
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

//
// Turn the passed instances into objects.
//
func CreateGameObjects(src M.InstanceMap,
) (ret GameObjects, err error,
) {
	ret = make(GameObjects)
	allProps := make(map[*M.ClassInfo]M.PropertySet)

MakeObjects:
	for _, info := range src {
		// create property sets for this instance's class
		class := info.Class()
		props, had := allProps[class]
		if !had {
			props = class.AllProperties()
			allProps[class] = props
		}
		// turn properties into tables:
		values, temps := NewRuntimeValues(), make(TemplatePool)
		values.set("Name", info.FullName())
		for propId, prop := range props {
			v, _ := info.ValueByName(prop.Name())
			switch val := v.(type) {
			case *M.EnumValue:
				choice := val.Choice()
				values.set(propId, choice)
				values.set(choice, true)
			case *M.NumValue:
				num, _ := val.Num()
				values.set(propId, num)
			case *M.TextValue:
				text, _ := val.Text()
				values.set(propId, text)
				// TBD: when to parse this? maybe options? here for early errors.
				if e := temps.New(propId.String(), text); e != nil {
					err = e
					break MakeObjects
				}
			case *M.RelativeValue:
				// no table enties
			default:
				err = fmt.Errorf("internal error: unknown property type %s:%T", propId, prop)
				break MakeObjects
			}
		}
		// creat the game obj
		gameobj := &GameObject{info, values, temps, E.NewDispatcher()}
		ret[info.Id()] = gameobj
	}
	return ret, err
}
