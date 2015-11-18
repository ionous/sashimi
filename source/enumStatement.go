package source

//
type EnumStatement struct {
	fields EnumFields
	source Code
}

//
type EnumFields struct {
	Class   string
	Name    string
	Choices []string
	//Expects []PropertyExpectation
}

//
func (ts EnumStatement) Source() Code {
	return ts.source
}

//
func (ts EnumStatement) Fields() EnumFields {
	return ts.fields
}

//
func (ts EnumFields) IndexOf(choice string) int {
	index := -1
	for i, v := range ts.Choices {
		if v == choice {
			index = i
			break
		}
	}
	return index
}
