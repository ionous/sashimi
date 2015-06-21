package source

type RelativeHint int

const (
	RelativeOne    = 1
	RelativeMany   = 2
	RelativeSource = 1 << 4
)

func (this RelativeHint) IsMany() bool {
	return this&RelativeMany != 0
}

func (this RelativeHint) IsReverse() bool {
	return this&RelativeSource == 0
}

//
type RelativeStatement struct {
	fields RelativeFields
	source Code
}

type RelativeFields struct {
	Class     string
	Property  string
	RelatesTo string
	Relation  string
	Hint      RelativeHint
}

//
func (this RelativeStatement) Fields() RelativeFields {
	return this.fields
}

//
func (this RelativeStatement) Source() Code {
	return this.source
}
