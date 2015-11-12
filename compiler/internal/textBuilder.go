package internal

import (
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/util/ident"
)

//
// NewTextBuilder returns an interface which can generate a text property
//
func NewTextBuilder(id ident.Id, name string) (IBuildProperty, error) {
	prop := M.TextProperty{id, name, false}
	return TextBuilder{prop}, nil
}

type TextBuilder struct {
	M.TextProperty
}

func (txt TextBuilder) BuildProperty() (M.IProperty, error) {
	return txt.TextProperty, nil
}

func (txt TextBuilder) SetProperty(ctx PropertyContext) (err error) {
	nilVal := ""
	return ctx.values.lockSet(ctx.inst, txt.Id, nilVal, ctx.value)
}
