package app

import (
	"github.com/ionous/sashimi/util/ident"
	"strings"
)

func jsonId(id ident.Id) string {
	return strings.Join(id.Split(), "-")
}
