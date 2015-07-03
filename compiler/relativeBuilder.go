package compiler

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
	S "github.com/ionous/sashimi/source"
)

type RelativeBuilder struct {
	cls       M.StringId
	id        M.StringId
	name      string
	src       S.Code
	fields    M.RelativeFields
	relations PendingRelations
}

func NewRelativeBuilder(
	relatives *RelativeFactory,
	cls, id M.StringId,
	name string,
	src S.Code,
	fields M.RelativeFields,
) (
	IBuildProperty,
	error,
) {
	return RelativeBuilder{cls, id, name, src, fields, relatives.relations}, nil
}

func (rel RelativeBuilder) BuildProperty() (ret M.IProperty, err error) {
	// FIX? the relation has split into two class relative properties
	// and now were merging them back together, verifying they match
	// it might have been better to keep the pair,
	// splitting them into their class halves at class creation time
	relationId := rel.fields.Relation
	relation := rel.relations[relationId]
	if e := relation.setRelative(rel.name, rel.fields); e != nil {
		panic(e)
		err = e
	} else {
		// write merged data back
		rel.relations[relationId] = relation
		// return the property
		ret = M.NewRelativeProperty(rel.fields)
	}
	return ret, err
}

func (rel RelativeBuilder) SetProperty(ctx PropertyContext) (err error) {
	nilVal := ""
	if otherName, okay := ctx.value.(string); !okay {
		err = SetValueMismatch(ctx.inst, rel.id, nilVal, ctx.value)
	} else {
		otherId := M.MakeStringId(otherName)
		if other, ok := ctx.refs[otherId]; !ok {
			err = M.InstanceNotFound(otherName)
		} else {
			err = rel._refSet(ctx.class, ctx.tables, ctx.inst, other)
		}
	}
	return err
}

// FIX: where and how to validate table.style:
// maybe in the table iteself since its already searching for duplicate pairs...
func (rel RelativeBuilder) _refSet(
	cls *M.ClassInfo,
	tables M.TableRelations,
	inst M.StringId,
	other *PartialInstance,
) (
	err error,
) {
	if !other.class.CompatibleWith(rel.fields.Relates) {
		// Claire.Pets value change 'gremlins' to 'Rocks'
		err = SetValueChanged(inst, rel.id, cls, rel.fields.Relates)
	} else if table, ok := tables.TableById(rel.fields.Relation); !ok {
		err = fmt.Errorf("internal error? couldn't find table for relation %+v", rel.fields.Relation)
	} else {
		src, dst := inst, other.id
		if rel.fields.IsRev {
			dst, src = src, dst
		}
		table.Add(src.String(), dst.String())
	}
	return err
}
