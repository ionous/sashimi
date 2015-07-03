package compiler

//
func newRelativeFactory(names NameScope) *RelativeFactory {
	return &RelativeFactory{names, make(PendingRelations)}
}

//
// RelativeFactory holds relative properties while building a picture of their potential linkage.
//
type RelativeFactory struct {
	NameScope
	relations PendingRelations
}
