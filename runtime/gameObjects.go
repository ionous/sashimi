package runtime

import (
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
		class := inst.Class
		props, had := allProps[class]
		if !had {
			props = class.AllProperties()
			allProps[class] = props
		}
		// turn properties into tables:
		gobj := &GameObject{inst.Id, inst.Class, make(TemplateValues), make(TemplatePool), tables}
		// FIX FIX FIX
		name := inst.Name
		if n, ok := inst.FindValue("printed name"); ok {
			if n, ok := n.(string); ok && n != "" {
				name = n
			}
		}
		gobj.setDirect("Name", name)

		for propId, prop := range props {
			val, _ := inst.Value(propId)
			if e := gobj.setValue(prop, val); e != nil {
				err = errutil.Append(err, e)
			}
		}
		ret[inst.Id] = gobj
	}
	return ret, err
}
