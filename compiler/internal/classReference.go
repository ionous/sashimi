package compiler

import (
	"fmt"
	M "github.com/ionous/sashimi/compiler/model"
	S "github.com/ionous/sashimi/source"
	"github.com/ionous/sashimi/util/errutil"
)

// ClassReferences stores a list of explicit class declarations
// (eg. for an PendingInstance)
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
	// loop over the list of class references
	var loc S.Code
	for _, ref := range refs.classes {
		if cls, ok := classes[ref.class.id]; !ok {
			err = errutil.Append(err, ClassNotFound(ref.class.String()))
		} else if class == nil {
			class, props = cls, ref.class.props
			loc = ref.source
		} else if class != cls {
			e := MultipleClassesNotSupported{class, cls, loc, ref.source}
			class, loc, err = cls, ref.source, errutil.Append(err, e)
		}
	}
	return class, props, err
}

type MultipleClassesNotSupported struct {
	one, two  *M.ClassInfo
	where, at S.Code
}

func (e MultipleClassesNotSupported) Error() string {
	err := fmt.Errorf("multiple classes not yet supported %v != %v", e.one, e.two)
	err = errutil.Append(err, SourceError(e.where, err))
	err = errutil.Append(err, SourceError(e.at, err))
	return err.Error()
}
