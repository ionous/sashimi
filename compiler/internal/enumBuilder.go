package compiler

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
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
		err = fmt.Errorf("internal error: expected enum contraints")
	} else {
		switch choice := ctx.value.(type) {
		case int:
			err = enum._enumSet(ctx.inst, choice, ctx.values, constraint)
		case string:
			choiceId := M.MakeStringId(choice)
			if idx, e := enum.choices.ChoiceToIndex(choiceId); e != nil {
				err = e
			} else {
				err = enum._enumSet(ctx.inst, idx, ctx.values, constraint)
			}
		default:
			var nilVal int = 0
			err = SetValueMismatch(ctx.inst, enum.id, nilVal, ctx.value)
		}
	}
	return err
}

func (enum EnumBuilder) _enumSet(inst ident.Id, choice int, values PendingValues, constraint *M.EnumConstraint) (err error) {
	if e := constraint.CheckIndex(choice); e != nil {
		err = e
	} else {
		var nilVal int = 0
		err = values.lockSet(inst, enum.id, nilVal, choice)
	}
	return err
}
