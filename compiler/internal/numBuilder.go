package internal

import (
	M "github.com/ionous/sashimi/compiler/xmodel"
	"github.com/ionous/sashimi/util/ident"
)

//
// NewNumBuilder returns an interface which can generate a number property
//
func NewNumBuilder(id ident.Id, name string, isMany bool) (IBuildProperty, error) {
	prop := M.NumProperty{id, name, isMany}
	return NumBuilder{prop}, nil
}

type NumBuilder struct {
	M.NumProperty
}

func (num NumBuilder) BuildProperty() (M.IProperty, error) {
	return num.NumProperty, nil
}

func (num NumBuilder) SetProperty(ctx PropertyContext) (err error) {
	nilVal := float64(0)
	switch val := ctx.value.(type) {
	case int:
		err = ctx.values.lockSet(ctx.inst, num.Id, nilVal, float64(val))
	case float32:
		err = ctx.values.lockSet(ctx.inst, num.Id, nilVal, float64(val))
	case float64: // note: go's own default number type is float64
		err = ctx.values.lockSet(ctx.inst, num.Id, nilVal, float64(val))
	default:
		err = SetValueMismatch(ctx.inst, num.Id, nilVal, val)
	}
	return err
}
