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
		state := []string{}
		content := []Dict{}
		props := make(Dict)
		// all properties is recursive through base classes
		for propId, prop := range cls.AllProperties() {
			// FIX: a faster value by id?
			propKey := propId.String()
			ivalue, _ := inst.ValueByName(prop.Name())
			switch val := ivalue.(type) {
			case *NumValue:
				num, _ := val.Num()
				props[propKey] = num

			case *EnumValue:
				choice := val.Choice()
				state = append(state, choice.String())

			case *TextValue:
				text, _ := val.Text()
				props[propKey] = text

			case *RelativeValue:
				recurse := map[string]bool{"Clothing": true, "Equipment": true, "Contents": true}
				if recurse[propKey] {
					for _, otherId := range val.List() {
						other := SerializeView(model, MakeStringId(otherId))
						content = append(content, other)
					}
				}
			}
		}
		ret = Dict{
			"id":      id.String(),
			"name":    inst.Name(),
			"cls":     cls.Id(),
			"state":   state,
			"props":   props,
			"content": content,
		}
	}
	return ret
}
