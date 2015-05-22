package script

type Parsing struct {
	phrases []string
}

func Matching(how string) Parsing {
	return Parsing{[]string{how}}
}

func (this Parsing) Or(words string) Parsing {
	this.phrases = append(this.phrases, words)
	return this
}
