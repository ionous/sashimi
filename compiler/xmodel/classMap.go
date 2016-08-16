package xmodel

import "github.com/ionous/sashimi/util/ident"

//
type ClassMap map[ident.Id]*ClassInfo

// helper to generate an escaped string and an error,
func (cm ClassMap) FindClass(name string) (ret *ClassInfo, okay bool) {
	id := MakeStringId(name)
	ret, okay = cm[id]
	return ret, okay
}

// FIX could be made faster
func (cm ClassMap) FindClassBySingular(name string) (ret *ClassInfo, okay bool) {
	for _, cls := range cm {
		if cls.Singular == name {
			ret, okay = cls, true
			break
		}
	}
	return ret, okay
}
