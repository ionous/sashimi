package runtime

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/util/errutil"
)

//
// Turn the passed instances into objects.
//
func CreateGameObjects(
	src M.InstanceMap,
	tables M.TableRelations,
) (
	ret GameObjects,
	err error,
) {
	ret = make(GameObjects)
	allProps := make(map[*M.ClassInfo]M.PropertySet)

	for _, inst := range src {
		// create property sets for this instance's class
		class := inst.Class()
		props, had := allProps[class]
		if !had {
			props = class.AllProperties()
			allProps[class] = props
		}
		// turn properties into tables:
		values := NewRuntimeValues(tables)

		values.setDirect("Name", inst.Name()) // full name?
		for propId, prop := range props {
			val, _ := inst.Value(propId)
			switch prop := prop.(type) {
			case *M.EnumProperty:
				if choice, e := prop.IndexToChoice(val.(int)); e != nil {
					err = errutil.Append(err, e)
				} else {
					values.setDirect(propId, choice)
					values.setDirect(choice, true)
				}
			case *M.NumProperty:
				values.setDirect(propId, val)
			case *M.PointerProperty:
				values.setDirect(propId, val)
			case *M.RelativeProperty:
				if table, ok := tables.TableById(prop.Relation()); !ok {
					e := fmt.Errorf("couldn't find table", prop.Relation())
					err = errutil.Append(err, e)
				} else {
					rel := RelativeValue{inst, prop, table}
					values.setDirect(propId, rel)
				}

			case *M.TextProperty:
				values.setDirect(propId, val)
				// TBD: when to parse this? maybe options? here for early errors.
				str := val.(string)
				if e := values.temps.New(propId.String(), str); e != nil {
					err = errutil.Append(err, e)
				}
			default:
				e := fmt.Errorf("internal error: unknown property type %s:%T", propId, prop)
				err = errutil.Append(err, e)
			}
		}
		// create the game obj
		gameobj := &GameObject{inst, values}
		ret[inst.Id()] = gameobj
	}
	return ret, err
}
