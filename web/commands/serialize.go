package commands

import (
	. "github.com/ionous/sashimi/model"
)

// FIX, FIX,FIX: this is wrong
// its serializing the original model, when actually we want to serialize the current values.

// FIX: cleanup. this is too much like PrintModel, like createObjects, and like what will be needed for tables.
// -- also how i originally stored the instance data ...
// we should have "templates" or procedures assigned as usually, needs work.
func SerializeView(model *Model, id StringId) (ret Dict) {
	if inst, ok := model.Instances[id]; ok {
		cls := inst.Class()
		states := []string{}
		values := make(Dict)
		// all properties is recursive through base classes
		for propId, prop := range cls.AllProperties() {
			// FIX: a faster value by id?
			propKey := propId.String()
			ivalue, _ := inst.ValueByName(prop.Name())
			switch val := ivalue.(type) {
			case *NumValue:
				num, _ := val.Num()
				values[propKey] = num

			case *EnumValue:
				choice := val.Choice()
				//values[propKey] = choice
				// values[choice] = true
				// states list, indexable values, both?
				states = append(states, choice.String())

			case *TextValue:
				text, _ := val.Text()
				values[propKey] = text

			case *RelativeValue:
				recurse := map[string]bool{"Clothing": true, "Equipment": true, "Contents": true}
				if recurse[propKey] {
					ar := []Dict{}
					for _, otherId := range val.List() {
						other := SerializeView(model, MakeStringId(otherId))
						ar = append(ar, other)
					}
					values[propKey] = ar
					// they all are...
					/*if relation.ToMany() {
						values[propKey] = ar
					} else if len(ar) > 0 {
						values[propKey] = ar[0]
					} else {
						values[propKey] = Dict{} // or null?
					}
					*/
				}
			}
		}
		values["_class"] = cls.Id() // note: plural
		values["_states"] = states
		ret = Dict{id.String(): values}
	}
	return ret
}
