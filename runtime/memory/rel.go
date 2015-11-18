package memory

import (
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/util/ident"
)

type relInfo struct {
	mdl *MemoryModel
	*M.RelationModel
}

func (r relInfo) GetId() ident.Id {
	return r.Id
}

// func (r *relInfo) RemoveRelative(src, dst ident.Id) {
// 	table := p.mdl.getTable(rel.Relation)
// 	isRev := r.getRelative(src).IsRev
// 	table.Remove(func(x, y ident.Id) bool {
// 		return (!isRev && dst == x) || (isRev && src == x)
// 	})
// }

// // FIX: im not a huge fan of the property search, its only needed for the auto-inversion:
// // FIX: im not a huge fan of the auto-inversion, can this be solved from the client side?
// func (r *relInfo) getRelative(src ident.Id) (srcProp *M.PropertyModel) {
// 	if src == r._relative.src {
// 		srcProp = r._relative.prop
// 	} else {
// 		inst := r.mdl.Instances[src]
// 		inst.Class.Visit(func(cls *M.ClassModel) (finished bool) {
// 			for _, p := range cls.Properties {
// 				if rel, ok := p.(M.PropertyModel); ok {
// 					if rel.Relation == r.Id {
// 						srcProp = &rel
// 						finished = true
// 						break
// 					}
// 				}
// 			}
// 			return
// 		})
// 		r._relative.src, r._relative.prop = src, srcProp
// 	}
// 	// outside the cache test to handle empty src id
// 	if srcProp == nil {
// 		panic(fmt.Sprintf("'%v' is not related by '%v'", src, r.Relation))
// 	}
// 	return
// }
