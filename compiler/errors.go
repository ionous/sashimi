package compiler

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/util/errutil"
)

//
// ClassNotFound
//
func ClassNotFound(class string) error {
	return errutil.Func(func() string {
		return fmt.Sprintf("class '%s' not found", class)
	})
}

//
// ClassNotFound
//
func EnumMultiplySpecified(class M.StringId, enum M.StringId) error {
	return errutil.Func(func() string {
		return fmt.Sprintf("enum %s.%s specified more than once", class, enum)
	})
}

//
// PropertyNotFound
//
func PropertyNotFound(class M.StringId, prop string) error {
	return errutil.Func(func() string {
		return fmt.Sprintf("property '%s.%s' not found", class, prop)
	})
}

//
// SetValueChanged
//
func SetValueChanged(inst, prop M.StringId, curr, want interface{}) error {
	return errutil.Func(func() string {
		return fmt.Sprintf("%s.%s value change '%v' to '%v'", inst, prop, curr, want)
	})
}

//
// SetValueMismatch
//
func SetValueMismatch(inst, prop M.StringId, want, got interface{}) error {
	return errutil.Func(func() string {
		return fmt.Sprintf("%s.%s expected value of %T got %T", inst, prop, want, got)
	})
}

//
// UnknownPropertyError
//
func UnknownPropertyError(cls M.StringId, name string) error {
	return errutil.Func(func() string {
		return fmt.Sprintf("unhandled property %s.%s.", cls, name)
	})
}
