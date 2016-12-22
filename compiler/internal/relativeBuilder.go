package internal

import (
	"github.com/ionous/sashimi/compiler/model/table"
	M "github.com/ionous/sashimi/compiler/xmodel"
	"github.com/ionous/sashimi/source/types"
	"github.com/ionous/sashimi/util/errutil"
	"github.com/ionous/sashimi/util/ident"
	"github.com/ionous/sashimi/util/sbuf"
)

type RelativeBuilder struct {
	cls       ident.Id
	id        ident.Id
	name      string
	src       types.Code
	fields    M.RelativeProperty
	relations PendingRelations
}

func NewRelativeBuilder(
	relatives *RelativeFactory,
	cls, id ident.Id,
	name string,
	src types.Code,
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
		err = errutil.New("couldnt find table", rel.fields.Relation)
	} else if otherName, okay := ctx.value.(string); !okay {
		err = errutil.New("relative builder", ctx.inst, rel.id, "invalid type", sbuf.Type{ctx.value}, ctx.value)
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
	table *table.Table,
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
