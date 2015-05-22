package script

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
func (this KnownAsFragment) And(name string) KnownAsFragment {
	this.names = append(this.names, name)
	return this
}

type KnownAsFragment struct {
	origin Origin
	names  []string
}

//
func (this KnownAsFragment) MakeStatement(b SubjectBlock) error {
	return b.NewAlias(b.subject, this.names)
}
