package internal

import (
	"github.com/ionous/mars/rt"
	M "github.com/ionous/sashimi/compiler/xmodel"
	"github.com/ionous/sashimi/util/errutil"
	"github.com/ionous/sashimi/util/ident"
)

//
// NewEnumBuilder returns an interface which can generate an enumerated property
//
func NewEnumBuilder(
	id ident.Id,
	name string,
	choices []string,
) (ret IBuildProperty, err error,
) {
	if enum, e := M.CheckedEnumeration(choices); e != nil {
		err = e
	} else {
		prop := M.EnumProperty{id, name, enum}
		ret = EnumBuilder{id, enum, prop}
	}
	return ret, err
}

type EnumBuilder struct {
	id      ident.Id
	choices M.Enumeration
	prop    M.EnumProperty
}

func (enum EnumBuilder) BuildProperty() (M.IProperty, error) {
	return enum.prop, nil
}

func (enum EnumBuilder) SetProperty(ctx PropertyContext) (err error) {

	var constraints M.IConstrain
	if c, ok := ctx.class.Constraints.ConstraintById(enum.id); ok {
		constraints = c
	} else {
		constraints = M.NewConstraint(enum.choices)
	}

	if constraint, ok := constraints.(*M.EnumConstraint); !ok {
		err = errutil.New("runtime error: expected enum contraints")
	} else {
		switch choice := ctx.value.(type) {
		case int:
			if c, e := constraint.IndexToChoice(choice); e != nil {
				err = e
				break
			} else {
				choiceId := rt.State(c)
				err = enum._enumSet(ctx.inst, choiceId, ctx.values, constraint)
			}
		case string:
			choiceId := rt.State(M.MakeStringId(choice))
			err = enum._enumSet(ctx.inst, choiceId, ctx.values, constraint)
		default:
			var nilVal int = 0
			err = SetValueMismatch("enum", ctx.inst, enum.id, nilVal, ctx.value)
		}
	}
	return err
}

func (enum EnumBuilder) _enumSet(inst ident.Id, choice rt.State, values PendingValues, constraint *M.EnumConstraint) (err error) {
	if e := constraint.CheckChoice(choice.Id()); e != nil {
		err = e
	} else {
		nilVal := (*rt.StateEval)(nil)
		err = values.lockSet(inst, enum.id, nilVal, choice)
	}
	return err
}
