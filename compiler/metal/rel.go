package metal

import (
	M "github.com/ionous/sashimi/compiler/model"
	"github.com/ionous/sashimi/util/ident"
)

type relInfo struct {
	mdl *Metal
	*M.RelationModel
}

func (r relInfo) GetId() ident.Id {
	return r.Id
}
