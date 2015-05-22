package model

import (
	"fmt"
)

//
type ClassMap map[StringId]*ClassInfo

// helper to generate an escaped string and an error,
func (this ClassMap) FindClass(name string) (ret *ClassInfo, okay bool) {
	id := MakeStringId(name)
	ret, okay = this[id]
	return ret, okay
}

// FIX could be made faster, also should return ok not error
func (this ClassMap) FindClassBySingular(name string) (ret *ClassInfo, err error) {
	for _, v := range this {
		if v.singular == name {
			ret = v
			break
		}
	}
	if ret == nil {
		err = fmt.Errorf("class not found `%s`", name)
	}
	return ret, err
}
