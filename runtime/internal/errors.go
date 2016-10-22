package internal

import (
	"github.com/ionous/sashimi/util/errutil"
)

func InstanceNotFound(name string) error {
	return errutil.New("instance not found", name)
}
