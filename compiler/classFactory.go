package compiler

import (
	"bitbucket.org/pkg/inflect"
	"fmt"
	M "github.com/ionous/sashimi/model"
	S "github.com/ionous/sashimi/source"
)

type PendingClasses map[M.StringId]*PendingClass // ptr for presense detection
type SingleToPlural map[string]M.StringId        // it kind of makes senses its single to plural id

type ClassFactory struct {
	allNames       NameSource
	relatives      *RelativeFactory
	pending        PendingClasses
	singleToPlural SingleToPlural
}

//
//
//
func newClassFactory(names NameSource, rel *RelativeFactory) *ClassFactory {
	res := &ClassFactory{names, rel, make(PendingClasses), make(SingleToPlural)}
	res.addClassRef(nil, "kinds", "kind")
	return res
}

//
// given the passed plural name, find the previously registered class
//
func (this *ClassFactory) findBySingularName(singular string,
) (*PendingClass, bool) {
	id := this.singleToPlural[singular]
	ret, okay := this.pending[id]
	return ret, okay
}

//
// given the passed plural name, find the previously registered class
//
func (this *ClassFactory) findByPluralName(plural string,
) (*PendingClass, bool) {
	id := M.MakeStringId(plural)
	ret, okay := this.pending[id]
	return ret, okay
}

//
//
//
func (this *ClassFactory) findByRelativeName(kind string, hint S.RelativeHint,
) (class *PendingClass, pluralized bool, err error) {
	switch hint & ^S.RelativeSource {
	case S.RelativeMany:
		if cls, ok := this.findByPluralName(kind); !ok {
			err = ClassNotFound(kind)
		} else {
			class, pluralized = cls, true
		}
	case S.RelativeOne:
		if cls, ok := this.findBySingularName(kind); !ok {
			err = ClassNotFound(kind)
		} else {
			class, pluralized = cls, false
		}
	default:
		if cls, ok := this.findByPluralName(kind); ok {
			class, pluralized = cls, true
		} else if cls, ok := this.findBySingularName(kind); ok {
			class, pluralized = cls, false
		} else {
			err = ClassNotFound(kind)
		}
	}
	return class, pluralized, err
}

//
func (this *ClassFactory) makeClasses(relatives *RelativeFactory) (
	classes M.ClassMap, err error,
) {
	cr := newResults(this.pending, relatives)
	return cr.finalizeClasses()
}

//
//
//
func (this *ClassFactory) addClassRef(parent *PendingClass, plural, single string,
) (class *PendingClass, err error,
) {
	// FIX: sanity check singular?
	if singular, e := this._addOptions(plural, single); e != nil {
		err = e
	} else if id, e := this.allNames.addName(nil, plural, "class", ""); e != nil {
		err = e
	} else if class = this.pending[id]; class != nil {
		// FIX? ratchet the class down?
		if class.parent != parent {
			err = fmt.Errorf("conflicting `%v` parent class `%v` respecified as `%v`",
				plural, class.parent, parent)
		}
	} else {
		var parentProps *PropertyBuilders
		if parent != nil {
			parentProps = &parent.props
		}
		class = &PendingClass{
			this, parent, id, plural, singular,
			this.allNames.newScope(plural),
			NewProperties(parentProps),
			make(PendingRules, 0),
		}
		this.pending[id] = class
		this.singleToPlural[singular] = id
	}

	return class, err
}

//
// ex. name="rooms", value="room".
//
func (this *ClassFactory) _addOptions(plural, singular string) (string, error) {
	if singular == "" {
		singular = inflect.Singularize(plural)
	}
	// reserve `room` to mean `rooms`
	// we dont return the id -- if they meant a specific singular string, we want that
	// the id is just the internals of name vs. name collision
	_, err := this.allNames.addName(nil, singular, plural, "")
	return singular, err
}
