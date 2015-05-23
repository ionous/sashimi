package compiler

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
	S "github.com/ionous/sashimi/source"
)

type PendingInstance struct {
	name     string
	longName string
	classes  ClassReferences
}

//
type ClassReferences struct {
	classes []ClassReference
}

type ClassReference struct {
	class  *PendingClass
	source S.Code // location of first reference
}

//
func (this *ClassReferences) addClassReference(class ClassReference) {
	found := false
	for _, v := range this.classes {
		if v.class == class.class {
			found = true
			break
		}
	}
	if !found {
		this.classes = append(this.classes, class)
	}
}

//
// given the finalized class map determine the most leaf class this instance represents
func (this ClassReferences) resolveClass(classes M.ClassMap,
) (ret *M.ClassInfo, err error,
) {
	if len(this.classes) != 1 {
		err = fmt.Errorf("multiple classes not yet supported; `%s`", this.classes)
	} else {
		classref := this.classes[0]
		if cls, ok := classes.FindClass(classref.class.name); !ok {
			err = fmt.Errorf("couldn't find class %s", classref.class.name)
		} else {
			ret = cls
		}
	}
	return ret, err
}
