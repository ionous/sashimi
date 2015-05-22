package script

import S "github.com/ionous/sashimi/source"

//
// Interface for generating statements.
//
type IFragment interface {
	MakeStatement(SubjectBlock) error
}

//
//
//
type SubjectBlock struct {
	subject    string
	theKeyword string
	*S.BuildingBlocks
}

//
// Respecify the interface as a function, for closures.
// ex.
// func SomeStatment(...) IFragment {
// return FunctionFragment{func(b SubjectBlock) error {return nil}}
// }
//
type FunctionFragment struct {
	call func(SubjectBlock) error
}

func (this FunctionFragment) MakeStatement(b SubjectBlock) error {
	return this.call(b)
}
