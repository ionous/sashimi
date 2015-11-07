package runtime

import (
	"github.com/ionous/sashimi/model/table"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/util/errutil"
)

// Turn the passed instances into objects.
func CreateGameObjects(
	mdl api.Model,
	tables table.Tables,
) (
	ret GameObjects,
	err error,
) {
	ret = make(GameObjects)
	printedName := MakeStringId("printed name")

	for i := 0; i < mdl.NumInstance(); i++ {
		inst := mdl.InstanceNum(i)
		// turn properties into tables:
		if gobj, e := NewGameObject(mdl, inst.GetId(), inst.GetParentClass(), inst, tables); e != nil {
			err = errutil.Append(err, e)
		} else {
			// FIX FIX FIX
			name := inst.GetOriginalName()
			if p, ok := inst.GetProperty(printedName); ok {
				if txt := p.GetValue().GetText(); txt != "" {
					name = txt
				}
			}
			gobj.setDirect(MakeStringId("name"), name)
			ret[inst.GetId()] = gobj
		}

	}
	return ret, err
}
