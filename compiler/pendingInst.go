package compiler

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
	S "github.com/ionous/sashimi/source"
	"github.com/ionous/sashimi/util/errutil"
)

type PendingInstance struct {
	id       M.StringId
	name     string
	longName string
	classes  ClassReferences
}

//
type ClassReferences struct {
	classes []ClassReference
}

type ClassReference struct {
	class  M.StringId
	source S.Code // location of first reference
}

//
func (this *ClassReferences) addClassReference(class M.StringId, source S.Code) {
	ref := ClassReference{class, source}
	this.classes = append(this.classes, ref)
}

//
// given the finalized class map determine the most leaf class this instance represents
func (this ClassReferences) resolveClass(classes M.ClassMap,
) (ret *M.ClassInfo, err error,
) {
	for _, ref := range this.classes {

		if cls, ok := classes[ref.class]; !ok {
			err = errutil.Append(err, ClassNotFound(ref.class.String()))
		} else if ret == nil {
			ret = cls
		} else if ret != cls {
			err = fmt.Errorf("multiple classes not yet supported; `%s`", this.classes)
		}
	}
	return ret, err
}
