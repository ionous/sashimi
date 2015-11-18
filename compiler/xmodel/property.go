package xmodel

import (
	"github.com/ionous/sashimi/util/ident"
	"strings"
)

// IProperty represents a sashimi type.
type IProperty interface {
	GetId() ident.Id
	GetName() string
}

type PropertySet map[ident.Id]IProperty

func (props PropertySet) GetPropertyByName(name string) (ret IProperty, okay bool) {
	for _, v := range props {
		if strings.EqualFold(name, v.GetName()) {
			ret, okay = v, true
			break
		}
	}
	return
}
