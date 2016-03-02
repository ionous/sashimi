package internal

import (
	"fmt"
	"github.com/ionous/sashimi/util/errutil"
)

func InstanceNotFound(name string) error {
	return errutil.Func(func() string {
		return fmt.Sprintf("instance not found `%s`", name)
	})
}
