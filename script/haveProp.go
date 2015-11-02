package script

import (
	S "github.com/ionous/sashimi/source"
	"strings"
)

//
// Statement to add a property to all instances of a class.
// FIX: currently, for relations, you must use HaveOne or HaveMany.
//
func Have(name string, kind string) ClassPropertyFragment {
	return ClassPropertyFragment{NewOrigin(2), name, kind}
}

//
// Establish a one-to-one, or one-to-many relation.
//
func HaveOne(name string, kind string) ClassRelationFragment {
	return ClassRelationFragment{src: ClassRelativeFragment{NewOrigin(2), name, kind, S.RelativeOne}}
}

//
// Establish a many-to-one relation.
//
func HaveMany(name string, kind string) ClassRelationFragment {
	return ClassRelationFragment{src: ClassRelativeFragment{NewOrigin(2), name, kind, S.RelativeMany}}
}

//
// pivot to add a reciprocal kind property relation
//
func (this ClassRelationFragment) Implying(kind string, dst ClassRelationFragment) IFragment {
	// NOTE: we can't test that the implied class matches the original have class
	// because the plurals might not match -- we rely on the compiler to detect mismatches
	return ClassRelationFragment{this.src, kind, dst.src}
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
func (this ClassPropertyFragment) MakeStatement(b SubjectBlock) error {
	fields := S.PropertyFields{b.subject, this.name, this.kind}
	return b.NewProperty(fields, this.origin.Code())
}

//
func (this ClassRelationFragment) MakeStatement(b SubjectBlock) (err error) {
	src, dst := this.src, this.dst
	// uses the subject, ex. gremlins, and the field, ex. pets: gremlins-pets-relation
	via := strings.Join([]string{b.subject, src.name, "relation"}, "-")

	srel := S.RelativeProperty{b.subject, src.name, src.kind, via, src.hint | S.RelativeSource}
	if e := b.NewRelative(srel, this.src.origin.Code()); e != nil {
		err = e
	} else if this.reverseClass != "" {
		drel := S.RelativeProperty{this.reverseClass, dst.name, dst.kind, via, dst.hint}
		err = b.NewRelative(drel, this.dst.origin.Code())
	}
	return err
}
