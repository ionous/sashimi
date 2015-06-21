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
	pending        PendingClasses
	singleToPlural SingleToPlural
}

//
//
//
func newClassFactory(names NameSource) *ClassFactory {
	res := &ClassFactory{names, make(PendingClasses), make(SingleToPlural)}
	res.addClassRef(nil, "kinds", S.Options{"singular name": "kind"})
	return res
}

//
// given the passed plural name, find the previously registered class
//
func (this *ClassFactory) findBySingularName(singular string,
) (class *PendingClass, err error) {
	plural := this.singleToPlural[singular]
	return this.findById(plural)
}

//
// given the passed plural name, find the previously registered class
//
func (this *ClassFactory) findByPluralName(plural string,
) (class *PendingClass, err error) {
	id := M.MakeStringId(plural)
	return this.findById(id)
}

//
//
//
func (this *ClassFactory) findByRelativeName(kind string, hint S.RelativeHint,
) (class *PendingClass, pluralized bool, err error) {
	switch hint & ^S.RelativeSource {
	case S.RelativeMany:
		if cls, e := this.findByPluralName(kind); e != nil {
			err = e
		} else {
			class, pluralized = cls, true
		}
	case S.RelativeOne:
		if cls, e := this.findBySingularName(kind); e != nil {
			err = e
		} else {
			class, pluralized = cls, false
		}
	default:
		if cls, e := this.findByPluralName(kind); e == nil {
			class, pluralized = cls, true
		} else if cls, e := this.findBySingularName(kind); e == nil {
			class, pluralized = cls, false
		} else {
			err = e
		}
	}
	return class, pluralized, err
}

//
func (this *ClassFactory) findById(id M.StringId) (class *PendingClass, err error) {
	class, ok := this.pending[id]
	if !ok {
		err = fmt.Errorf("couldn't find class named `%s`", id)
	}
	return class, err
}

//
func (this *ClassFactory) makeClasses(relatives *RelativeFactory) (
	classes M.ClassMap,
	err error,
) {
	cr := newResults(this.pending, relatives)
	return cr.finalizeClasses()
}

//
//
//
func (this *ClassFactory) addClassRef(parent *PendingClass, plural string, options S.Options,
) (class *PendingClass, err error,
) {
	// FIX: sanity check singular?
	if singular, e := this._addOptions(plural, options); e != nil {
		err = e
	} else {
		// validate the new class plural
		if id, e := this.allNames.addName(nil, plural, "class", ""); e != nil {
			err = e
		} else {
			// ensure the class exists
			class = this.pending[id]
			if class != nil {
				// FIX? ratchet the class down?
				if class.parent != parent {
					err = fmt.Errorf("conflicting `%v` parent class `%v` respecified as `%v`",
						plural, class.parent, parent)
				}
			} else {
				class = &PendingClass{
					this, parent, id, plural, singular,
					this.allNames.newScope(plural),
					make(PendingProperties),
					make(PendingRules, 0), //PendingRules{},
					make(PendingRelatives),
				}
				this.pending[id] = class
				this.singleToPlural[singular] = id
			}
		}
	}
	return class, err
}

//
// ex. name="rooms", value="room".
//
func (this *ClassFactory) _addOptions(plural string, options S.Options,
) (singular string, err error,
) {
	singular = options["singular name"]
	if singular == "" {
		singular = inflect.Singularize(plural)
	}
	// reserve `room` to mean `rooms`
	// we dont return the id -- if they meant a specific singular string, we want that
	// the id is just the internals of name vs. name collision
	_, err = this.allNames.addName(nil, singular, plural, "")
	return singular, err
}
