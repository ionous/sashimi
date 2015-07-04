package runtime

import (
	"fmt"
	E "github.com/ionous/sashimi/event"
	M "github.com/ionous/sashimi/model"
)

//
// Turn the passed instances into objects.
//
func CreateGameObjects(src M.InstanceMap,
) (ret GameObjects, err error,
) {
	ret = make(GameObjects)
	allProps := make(map[*M.ClassInfo]M.PropertySet)

MakeObjects:
	for _, inst := range src {
		// create property sets for this instance's class
		class := inst.Class()
		props, had := allProps[class]
		if !had {
			props = class.AllProperties()
			allProps[class] = props
		}
		// turn properties into tables:
		values, temps := NewRuntimeValues(), make(TemplatePool)
		values.setDirect("Name", inst.FullName())
		for propId, prop := range props {
			val := inst.PropertyValue(prop)
			switch prop := prop.(type) {
			case *M.PointerProperty:
				values.setDirect(propId, val)
			case *M.NumProperty:
				values.setDirect(propId, val)
			case *M.TextProperty:
				values.setDirect(propId, val)
				// TBD: when to parse this? maybe options? here for early errors.
				str := val.(string)
				if e := temps.New(propId.String(), str); e != nil {
					err = e
					break MakeObjects
				}
			case *M.EnumProperty:
				choice := val.(M.StringId)
				values.setDirect(propId, choice)
				values.setDirect(choice, true)
			case *M.RelativeProperty:
				// no table enties
			default:
				err = fmt.Errorf("internal error: unknown property type %s:%T", propId, prop)
				break MakeObjects
			}
		}
		// creat the game obj
		gameobj := &GameObject{inst, values, temps, E.NewDispatcher()}
		ret[inst.Id()] = gameobj
	}
	return ret, err
}
