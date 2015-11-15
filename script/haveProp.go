package script

import (
	S "github.com/ionous/sashimi/source"
	"strings"
)

//
// Have statement to add a property to all instances of a class.
// FIX: currently, for relations, you must use HaveOne or HaveMany.
//
func Have(name string, kind string) ClassPropertyFragment {
	return ClassPropertyFragment{NewOrigin(2), name, kind}
}

// HaveOne establishes a one-to-one, or one-to-many relation.
func HaveOne(name string, kind string) ClassRelationFragment {
	return ClassRelationFragment{src: ClassRelativeFragment{NewOrigin(2), name, kind, S.RelativeOne}}
}

// HaveMany establishs a many-to-one relation.
func HaveMany(name string, kind string) ClassRelationFragment {
	return ClassRelationFragment{src: ClassRelativeFragment{NewOrigin(2), name, kind, S.RelativeMany}}
}

// Pivot to add a reciprocal kind property relation
func (frag ClassRelationFragment) Implying(kind string, dst ClassRelationFragment) IFragment {
	// NOTE: we can't test that the implied class matches the original have class
	// because the plurals might not match -- we rely on the compiler to detect mismatches
	return ClassRelationFragment{frag.src, kind, dst.src}
}

//
type ClassPropertyFragment struct {
	origin Origin
	name   string // property,field name
	kind   string // property kind: primitive or user class
}

type ClassRelativeFragment struct {
	origin Origin
	name   string // property,field name
	kind   string // property kind: primitive or user class
	hint   S.RelativeHint
}

//
type ClassRelationFragment struct {
	src          ClassRelativeFragment
	reverseClass string // the reverse subject from implying(reverse)
	dst          ClassRelativeFragment
}

//
func (frag ClassPropertyFragment) MakeStatement(b SubjectBlock) error {
	// FIX? kind of hacky, maybe should be HaveMany() instead
	list, isMany := " list", false
	if i := strings.Index(frag.kind, list); i > 0 {
		if i+len(list) == len(frag.kind) {
			frag.kind = frag.kind[:i]
			isMany = true
		}
	}
	fields := S.PropertyFields{b.subject, frag.name, frag.kind, isMany}
	return b.NewProperty(fields, frag.origin.Code())
}

//
func (frag ClassRelationFragment) MakeStatement(b SubjectBlock) (err error) {
	src, dst := frag.src, frag.dst
	// uses the subject, ex. gremlins, and the field, ex. pets: gremlins-pets-relation
	via := strings.Join([]string{b.subject, src.name, "relation"}, "-")

	srel := S.RelativeProperty{b.subject, src.name, src.kind, via, src.hint | S.RelativeSource}
	if e := b.NewRelative(srel, frag.src.origin.Code()); e != nil {
		err = e
	} else if frag.reverseClass != "" {
		drel := S.RelativeProperty{frag.reverseClass, dst.name, dst.kind, via, dst.hint}
		err = b.NewRelative(drel, frag.dst.origin.Code())
	}
	return err
}
