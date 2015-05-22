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
	Expects []PropertyExpectation
}

//
func (this EnumStatement) Source() Code {
	return this.source
}

//
func (this EnumStatement) Fields() EnumFields {
	return this.fields
}

//
func (this EnumFields) IndexOf(choice string) int {
	index := -1
	for i, v := range this.Choices {
		if v == choice {
			index = i
			break
		}
	}
	return index
}
