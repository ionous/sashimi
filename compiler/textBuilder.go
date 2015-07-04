package compiler

import (
	M "github.com/ionous/sashimi/model"
)

//
// NewTextBuilder returns an interface which can generate a text property
//
func NewTextBuilder(id M.StringId, name string) (IBuildProperty, error) {
	prop := M.NewTextProperty(id, name)
	return TextBuilder{prop}, nil
}

type TextBuilder struct {
	prop *M.TextProperty
}

func (txt TextBuilder) BuildProperty() (M.IProperty, error) {
	return txt.prop, nil
}

func (txt TextBuilder) SetProperty(ctx PropertyContext) (err error) {
	nilVal := ""
	return ctx.values.lockSet(ctx.inst, txt.prop.Id(), nilVal, ctx.value)
}
