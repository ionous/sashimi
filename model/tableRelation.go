package model

import (
	"github.com/ionous/sashimi/model/table"
	"github.com/ionous/sashimi/util/ident"
)

//
type TableRelations struct {
	relation RelationMap // FIX: a little weird this is here when it's also in the model itself
	tables   map[ident.Id]*TableRelation
}

//
type TableRelation struct {
	*Relation
	*table.Table
}

//
func NewTableRelations(rels RelationMap) TableRelations {
	tables := make(map[ident.Id]*TableRelation)
	for id, rel := range rels {
		entry := &TableRelation{&rel, table.NewTable()}
		tables[id] = entry
	}
	return TableRelations{rels, tables}
}

//
func (this TableRelations) TableById(id ident.Id) (ret *TableRelation, ok bool) {
	ret, ok = this.tables[id]
	return ret, ok
}
