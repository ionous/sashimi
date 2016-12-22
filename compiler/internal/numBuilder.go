package internal

import (
	"github.com/ionous/mars/rt"
	M "github.com/ionous/sashimi/compiler/xmodel"
	"github.com/ionous/sashimi/util/errutil"
	"github.com/ionous/sashimi/util/ident"
	"reflect"
)

// NewNumBuilder returns an interface which can generate a number property.
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
	if !num.IsMany {
		nilVal := (*rt.NumberEval)(nil)
		switch val := ctx.value.(type) {
		case int:
			err = ctx.values.lockSet(ctx.inst, num.Id, nilVal, rt.Number{float64(val)})
		case float32:
			err = ctx.values.lockSet(ctx.inst, num.Id, nilVal, rt.Number{float64(val)})
		case float64: // note: go's own default number type is float64
			err = ctx.values.lockSet(ctx.inst, num.Id, nilVal, rt.Number{val})
		case rt.NumberEval:
			err = ctx.values.lockSet(ctx.inst, num.Id, nilVal, val)
		default:
			err = errutil.New("NumBuilder: unexpected type", ctx.inst, num.Id, reflect.TypeOf(val))
		}
	} else {
		nilVal := (*rt.NumListEval)(nil)
		switch val := ctx.value.(type) {
		case rt.NumListEval:
			err = ctx.values.lockSet(ctx.inst, num.Id, nilVal, val)
		default:
			err = errutil.New("NumBuilder: unexpected list type", ctx.inst, num.Id, val)
		}
	}
	return err
}
