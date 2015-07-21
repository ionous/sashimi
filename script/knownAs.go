package script

import (
	S "github.com/ionous/sashimi/source"
)

//
// Add an alias for the current subject.
// ex. The("cabinet", IsKnownAs("armoire").And("..."))
//
func IsKnownAs(name string) KnownAsFragment {
	return KnownAsFragment{NewOrigin(1), append([]string{name})}
}

//
// Add additional aliases for the current subject.
//
func (frag KnownAsFragment) And(name string) KnownAsFragment {
	frag.names = append(frag.names, name)
	return frag
}

type KnownAsFragment struct {
	Origin
	names []string
}

//
func (frag KnownAsFragment) MakeStatement(b SubjectBlock) error {
	alias := S.AliasFields{b.subject, frag.names}
	return b.NewAlias(alias, frag.Code())
}
