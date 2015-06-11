package model

//
type ClassMap map[StringId]*ClassInfo

// helper to generate an escaped string and an error,
func (this ClassMap) FindClass(name string) (ret *ClassInfo, okay bool) {
	id := MakeStringId(name)
	ret, okay = this[id]
	return ret, okay
}

// FIX could be made faster
func (this ClassMap) FindClassBySingular(name string) (ret *ClassInfo, okay bool) {
	for _, cls := range this {
		if cls.singular == name {
			ret, okay = cls, true
			break
		}
	}
	return ret, okay
}
