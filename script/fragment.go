package script

import S "github.com/ionous/sashimi/source"

//
// Interface for generating statements.
//
type IFragment interface {
	MakeStatement(SubjectBlock) error
}

//
type SubjectBlock struct {
	*Script
	theKeyword string
	subject    string
	*S.BuildingBlocks
}

func (sb *SubjectBlock) Keyword() string {
	return sb.theKeyword
}

func (sb *SubjectBlock) Subject() string {
	return sb.subject
}

//
// NewFunctionFragment changes the passed function into an implementation of IFragment.
//
func NewFunctionFragment(cb functionFragmentCall) IFragment {
	return _FunctionFragment{cb}
}

type functionFragmentCall func(SubjectBlock) error

type _FunctionFragment struct {
	call functionFragmentCall
}

func (fcb _FunctionFragment) MakeStatement(b SubjectBlock) error {
	return fcb.call(b)
}
