package compiler

import (
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/util/ident"
)

//
// NewNumBuilder returns an interface which can generate a number property
//
func NewNumBuilder(id ident.Id, name string) (IBuildProperty, error) {
	prop := M.NumProperty{id, name}
	return NumBuilder{prop}, nil
}

type NumBuilder struct {
	M.NumProperty
}

func (num NumBuilder) BuildProperty() (M.IProperty, error) {
	return num.NumProperty, nil
}

func (num NumBuilder) SetProperty(ctx PropertyContext) (err error) {
	nilVal := float32(0)
	switch val := ctx.value.(type) {
	case int:
		err = ctx.values.lockSet(ctx.inst, num.Id, nilVal, float32(val))
	case float32:
		err = ctx.values.lockSet(ctx.inst, num.Id, nilVal, float32(val))
	case float64: // note: go's own default number type is float64
		err = ctx.values.lockSet(ctx.inst, num.Id, nilVal, float32(val))
	default:
		err = SetValueMismatch(ctx.inst, num.Id, nilVal, val)
	}
	return err
}
