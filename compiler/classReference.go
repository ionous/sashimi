package compiler

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
	S "github.com/ionous/sashimi/source"
	"github.com/ionous/sashimi/util/errutil"
)

//
type ClassReferences struct {
	classes []ClassReference
}

type ClassReference struct {
	class  *PendingClass
	source S.Code // location of first reference
}

//
func (refs *ClassReferences) addClassReference(class *PendingClass, source S.Code) {
	ref := ClassReference{class, source}
	refs.classes = append(refs.classes, ref)
}

//
// given the finalized class map determine the most leaf class refs instance represents
func (refs ClassReferences) resolveClass(classes M.ClassMap,
) (class *M.ClassInfo, props PropertyBuilders, err error,
) {
	for _, ref := range refs.classes {
		if cls, ok := classes[ref.class.id]; !ok {
			err = errutil.Append(err, ClassNotFound(ref.class.String()))
		} else if class == nil {
			class, props = cls, ref.class.props
		} else if class != cls {
			err = fmt.Errorf("multiple classes not yet supported; `%s`", refs.classes)
		}
	}
	return class, props, err
}
