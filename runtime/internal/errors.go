package internal

import (
	"fmt"
	"github.com/ionous/sashimi/util/errutil"
)

func TypeMismatch(name string, kind string) error {
	return errutil.Func(func() string {
		return fmt.Sprintf("type mismatch %s expected %s", name, kind)
	})
}

func InstanceNotFound(name string) error {
	return errutil.Func(func() string {
		return fmt.Sprintf("instance not found `%s`", name)
	})
}
