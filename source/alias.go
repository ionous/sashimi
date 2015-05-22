package source

//
type AliasStatement struct {
	key     string
	phrases []string
	source  Code
}

//
func (this AliasStatement) Key() string {
	return this.key
}

//
func (this AliasStatement) Phrases() []string {
	return this.phrases
}

//
func (this AliasStatement) Source() Code {
	return this.source
}
