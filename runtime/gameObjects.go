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
		values.setDirect("Name", info.FullName())
		for propId, prop := range props {
			v, _ := info.ValueByName(prop.Name())
			switch val := v.(type) {
			case *M.EnumValue:
				choice := val.Choice()
				values.setDirect(propId, choice)
				values.setDirect(choice, true)
			case *M.NumValue:
				num, _ := val.Num()
				values.setDirect(propId, num)
			case *M.TextValue:
				text, _ := val.Text()
				values.setDirect(propId, text)
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
