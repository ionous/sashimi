package source

//
type AliasStatement struct {
	key     string
	phrases []string
	source  Code
}

//
func (ts AliasStatement) Key() string {
	return ts.key
}

//
func (ts AliasStatement) Phrases() []string {
	return ts.phrases
}

//
func (ts AliasStatement) Source() Code {
	return ts.source
}
