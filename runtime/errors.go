package runtime

import (
	"fmt"
	"github.com/ionous/sashimi/util/errutil"
)

//
//
//
func TypeMismatch(name string, kind string) error {
	return errutil.Func(func() string {
		return fmt.Sprintf("type mismatch %s expected %s", name, kind)
	})
}

//
//
//
func NoSuchValue(owner string, choice string) error {
	return errutil.Func(func() string {
		return fmt.Sprintf("no such value '%s'.'%s'", owner, choice)
	})
}

//
// for unexpected runtime errors
//
func RuntimeError(err error) error {
	return err
}
