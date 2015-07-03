package compiler

import (
	M "github.com/ionous/sashimi/model"
)

func NewPointerBuilder(id M.StringId, name string, class M.StringId) IBuildProperty {
	prop := M.NewPointerProperty(id, name, class)
	return PointerBuilder{id, class, prop}
}

type PointerBuilder struct {
	id    M.StringId
	class M.StringId
	prop  *M.PointerProperty
}

func (ptr PointerBuilder) BuildProperty() (M.IProperty, error) {
	return ptr.prop, nil
}

func (ptr PointerBuilder) SetProperty(ctx PropertyContext) (err error) {
	nilVal := ""
	if otherName, okay := ctx.value.(string); !okay {
		err = SetValueMismatch(ctx.inst, ptr.id, nilVal, ctx.value)
	} else {
		otherId := M.MakeStringId(otherName)
		if _, ok := ctx.refs[otherId]; !ok {
			err = M.InstanceNotFound(otherName)
		} else {
			err = ctx.values.lockSet(ctx.inst, ptr.prop.Id(), nilVal, otherId)
		}
	}
	return err
}
