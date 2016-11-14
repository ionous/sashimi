package internal

import (
	"github.com/ionous/mars/rt"
	M "github.com/ionous/sashimi/compiler/xmodel"
	"github.com/ionous/sashimi/util/errutil"
	"github.com/ionous/sashimi/util/ident"
)

// NewTextBuilder returns an interface which can generate a text property.
func NewTextBuilder(id ident.Id, name string, isMany bool) (IBuildProperty, error) {
	prop := M.TextProperty{id, name, isMany}
	return TextBuilder{prop}, nil
}

type TextBuilder struct {
	M.TextProperty
}

func (txt TextBuilder) BuildProperty() (M.IProperty, error) {
	return txt.TextProperty, nil
}

func (txt TextBuilder) SetProperty(ctx PropertyContext) (err error) {
	if !txt.IsMany {
		nilVal := (*rt.TextEval)(nil)
		switch val := ctx.value.(type) {
		case string:
			err = ctx.values.lockSet(ctx.inst, txt.Id, nilVal, rt.Text(val))
		case rt.TextEval:
			err = ctx.values.lockSet(ctx.inst, txt.Id, nilVal, val)
		default:
			err = errutil.New("TextBuilder: unexpected type", ctx.inst, txt.Id, val)
		}
	} else {
		nilVal := (*rt.TextListEval)(nil)
		switch val := ctx.value.(type) {
		case rt.TextListEval:
			err = ctx.values.lockSet(ctx.inst, txt.Id, nilVal, val)
		default:
			err = errutil.New("TextBuilder: unexpected list type", ctx.inst, txt.Id, val)
		}
	}
	return
}
