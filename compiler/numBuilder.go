package compiler

import (
	M "github.com/ionous/sashimi/model"
)

//
// NewNumBuilder returns an interface which can generate a number property
//
func NewNumBuilder(id M.StringId, name string) (IBuildProperty, error) {
	prop := M.NewNumProperty(id, name)
	return NumBuilder{id, prop}, nil
}

type NumBuilder struct {
	id   M.StringId
	prop *M.NumProperty
}

func (num NumBuilder) BuildProperty() (M.IProperty, error) {
	return num.prop, nil
}

func (num NumBuilder) SetProperty(ctx PropertyContext) (err error) {
	nilVal := float32(0)
	switch val := ctx.value.(type) {
	case int:
		err = ctx.values.lockSet(ctx.inst, num.id, nilVal, float32(val))
	case float32:
		err = ctx.values.lockSet(ctx.inst, num.id, nilVal, float32(val))
	case float64: // note: go's own default number type is float64
		err = ctx.values.lockSet(ctx.inst, num.id, nilVal, float32(val))
	default:
		err = SetValueMismatch(ctx.inst, num.id, nilVal, val)
	}
	return err
}
