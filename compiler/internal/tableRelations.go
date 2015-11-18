package internal

import (
	M "github.com/ionous/sashimi/compiler/xmodel"
	"github.com/ionous/sashimi/model/table"
	"github.com/ionous/sashimi/util/ident"
)

type TableRelations struct {
	RelationMap M.RelationMap
	Tables      table.Tables
}

//
func NewTableRelations(rels M.RelationMap) TableRelations {
	tables := make(table.Tables)
	for id, _ := range rels {
		tables[id] = table.NewTable()
	}
	return TableRelations{rels, tables}
}

func (t TableRelations) TableById(id ident.Id) (ret *table.Table, ok bool) {
	ret, ok = t.Tables[id]
	return ret, ok
}
