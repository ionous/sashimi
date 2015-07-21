package compiler

import (
	"fmt"
	"github.com/ionous/sashimi/source"
	"github.com/ionous/sashimi/util/errutil"
	"github.com/ionous/sashimi/util/ident"
)

func SourceError(src source.Code, e error) error {
	return errutil.Func(func() string {
		return fmt.Sprintf("error! %s @ %s", e.Error(), src)
	})
}

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
func EnumMultiplySpecified(class ident.Id, enum ident.Id) error {
	return errutil.Func(func() string {
		return fmt.Sprintf("enum %s.%s specified more than once", class, enum)
	})
}

//
// PropertyNotFound
//
func PropertyNotFound(class ident.Id, prop string) error {
	return errutil.Func(func() string {
		return fmt.Sprintf("property '%s.%s' not found", class, prop)
	})
}

//
// SetValueChanged
//
func SetValueChanged(inst, prop ident.Id, curr, want interface{}) error {
	return errutil.Func(func() string {
		return fmt.Sprintf("%s.%s value change '%v' to '%v'", inst, prop, curr, want)
	})
}

//
// SetValueMismatch
//
func SetValueMismatch(inst, prop ident.Id, want, got interface{}) error {
	return errutil.Func(func() string {
		return fmt.Sprintf("%s.%s expected value of %T got %T", inst, prop, want, got)
	})
}

//
// UnknownPropertyError
//
func UnknownPropertyError(cls ident.Id, name string) error {
	return errutil.Func(func() string {
		return fmt.Sprintf("unhandled property %s.%s.", cls, name)
	})
}
