package script

import (
	S "github.com/ionous/sashimi/source"
	"strings"
)

//
// Statement to add a property to all instances of a class
// for relations: we start by saying how many of the target side we have
// and, imply controls whether this source side is many or (default) one
//
func Have(name string, kind string) ClassPropertyFragment {
	return ClassPropertyFragment{name, kind, S.RelativeGuess}
}

//
func HaveOne(name string, kind string) ClassPropertyFragment {
	return ClassPropertyFragment{name, kind, S.RelativeOne}
}

//
func HaveMany(name string, kind string) ClassPropertyFragment {
	return ClassPropertyFragment{name, kind, S.RelativeMany}
}

//
// pivot to add a reciprocal kind property relation
//
func (this ClassPropertyFragment) Implying(kind string, dst ClassPropertyFragment) (frag IFragment) {
	// NOTE: we can't test that the implied class matches the original have class
	// because the plurals might not match -- we rely on the compiler to detect mismatches
	return ClassRelationFragment{this, kind, dst}
}

//
//
//
type ClassPropertyFragment struct {
	name string // property,field name
	kind string // property kind: primitive or user class
	hint S.RelativeHint
}

//
type ClassRelationFragment struct {
	src          ClassPropertyFragment
	reverseClass string // the reverse subject from implying(reverse)
	dst          ClassPropertyFragment
}

//
func (this ClassPropertyFragment) MakeStatement(b SubjectBlock) error {
	fields := S.PropertyFields{b.subject, this.name, this.kind, this.hint | S.RelativeSource}
	return b.NewProperty(fields, "")
}

//
func (this ClassRelationFragment) MakeStatement(b SubjectBlock) (err error) {

	src, dst := this.src, this.dst
	// uses the subject, ex. gremlins, and the field, ex. pets
	via := strings.Join([]string{b.subject, src.name, "relation"}, "-")

	srel := S.RelativeFields{b.subject, src.name, src.kind, via, src.hint | S.RelativeSource}
	if e := b.NewRelative(srel, ""); e != nil {
		err = e
	} else {
		drel := S.RelativeFields{this.reverseClass, dst.name, dst.kind, via, dst.hint}
		err = b.NewRelative(drel, "")
	}
	return err
}
