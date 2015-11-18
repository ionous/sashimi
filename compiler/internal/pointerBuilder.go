package internal

import (
	M "github.com/ionous/sashimi/compiler/xmodel"
	"github.com/ionous/sashimi/util/ident"
)

// note: it probably doesnt make sense to allow a ratcheting down of cls
// such a ratcheting would *increasing* restrictiveness, not permissability.
// for example: "pointer","kind" could store "teddy bears",
// but changed to "pointer","adult white male" could only store "teddy roosevelt"
func NewPointerBuilder(id ident.Id, name string, class ident.Id, isMany bool) (IBuildProperty, error) {
	prop := M.PointerProperty{id, name, class, isMany}
	return PointerBuilder{prop}, nil
}

type PointerBuilder struct {
	M.PointerProperty
}

func (ptr PointerBuilder) BuildProperty() (M.IProperty, error) {
	return ptr.PointerProperty, nil
}

func (ptr PointerBuilder) SetProperty(ctx PropertyContext) (err error) {
	var nilVal ident.Id
	if otherName, okay := ctx.value.(string); !okay {
		err = SetValueMismatch(ctx.inst, ptr.Id, nilVal, ctx.value)
	} else {
		otherId := M.MakeStringId(otherName)
		if _, ok := ctx.refs[otherId]; !ok {
			err = M.InstanceNotFound(otherName)
		} else {
			err = ctx.values.lockSet(ctx.inst, ptr.Id, nilVal, otherId)
		}
	}
	return err
}
