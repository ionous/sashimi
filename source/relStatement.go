package source

type RelativeHint int

const (
	RelativeOne    = 1
	RelativeMany   = 2
	RelativeSource = 1 << 4
)

func (hint RelativeHint) IsMany() bool {
	return hint&RelativeMany != 0
}

func (hint RelativeHint) IsReverse() bool {
	return hint&RelativeSource == 0
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
func (ts RelativeStatement) Fields() RelativeFields {
	return ts.fields
}

//
func (ts RelativeStatement) Source() Code {
	return ts.source
}
