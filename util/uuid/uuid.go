//Package uuid isolates the crypto libraries from the rest of the utils
// TODO: an "id maker" interface implemented by various things?
package uuid

import (
	"github.com/ionous/sashimi/util/ident"
	"github.com/satori/go.uuid"
)

func MakeUniqueId() ident.Id {
	str := uuid.NewV4().String()
	return ident.Id("~" + str)
}
