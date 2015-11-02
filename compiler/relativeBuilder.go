package compiler

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
	S "github.com/ionous/sashimi/source"
	"github.com/ionous/sashimi/util/ident"
)

type RelativeBuilder struct {
	cls       ident.Id
	id        ident.Id
	name      string
	src       S.Code
	fields    M.RelativeProperty
	relations PendingRelations
}

func NewRelativeBuilder(
	relatives *RelativeFactory,
	cls, id ident.Id,
	name string,
	src S.Code,
	fields M.RelativeProperty,
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
	relative := PendingRelative{rel.cls, rel.fields}
	if e := relation.setRelative(rel.name, relative); e != nil {
		err = e
	} else {
		// write merged data back
		rel.relations[relationId] = relation
		// return the property
		ret = rel.fields
	}
	return ret, err
}

func (rel RelativeBuilder) SetProperty(ctx PropertyContext) (err error) {
	if table, ok := ctx.tables.TableById(rel.fields.Relation); !ok {
		err = fmt.Errorf("couldnt find table", rel.fields.Relation)
	} else if otherName, okay := ctx.value.(string); !okay {
		err = SetValueMismatch(ctx.inst, rel.id, "", ctx.value)
	} else {
		otherId := M.MakeStringId(otherName)
		if other, ok := ctx.refs[otherId]; !ok {
			err = M.InstanceNotFound(otherName)
		} else {
			err = rel._refSet(table, ctx.class, ctx.inst, other)
		}
	}
	return err
}

// FIX: where and how to validate table.style:
// maybe in the table iteself since its already searching for duplicate pairs...
func (rel RelativeBuilder) _refSet(
	table *M.TableRelation,
	cls *M.ClassInfo,
	inst ident.Id,
	other *PartialInstance,
) (
	err error,
) {
	if !other.class.CompatibleWith(rel.fields.Relates) {
		// Claire.Pets value change 'gremlins' to 'Rocks'
		err = SetValueChanged(inst, rel.id, cls, rel.fields.Relates)
	} else {
		src, dst := inst, other.id
		if rel.fields.IsRev {
			dst, src = src, dst
		}
		table.Add(src, dst)
	}
	return err
}
