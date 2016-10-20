package internal

import (
	"github.com/ionous/sashimi/source"
	"github.com/ionous/sashimi/util/errutil"
	"github.com/ionous/sashimi/util/ident"
)

func SourceError(src source.Code, err error) error {
	return errutil.New("source code error", src, err)
}

func ClassNotFound(class string) error {
	return errutil.New("class", class, "not found")
}

func EnumMultiplySpecified(class ident.Id, enum ident.Id) error {
	return errutil.New("class enum", class, enum, "specified more than once")
}

func PropertyNotFound(class ident.Id, prop string) error {
	return errutil.New("class property", class, prop, "property not found")
}

func SetValueChanged(inst, prop ident.Id, curr, want interface{}) error {
	return errutil.New("instance propeperty", inst, prop, "value change", curr, "to", want)
}

func SetValueMismatch(name, inst, prop ident.Id, want, got interface{}) error {
	return errutil.New("instance property", inst, prop, name, "expected value of", want, "got", got)
}

func UnknownPropertyError(class ident.Id, name string) error {
	return errutil.New("class property", class, name, "unknown error")
}
