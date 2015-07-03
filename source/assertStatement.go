package source

//
type AssertionStatement struct {
	owner   string  // base type or class
	called  string  // name of reference being asserted into existance
	options Options // ex. called
}

//
func (this AssertionStatement) Owner() string {
	return this.owner
}

//
// bare name without articles
func (this AssertionStatement) ShortName() string {
	return this.called
}

//
func (this AssertionStatement) FullName() string {
	long := this.options["long name"]
	if long == "" {
		long = this.called
	}
	return long
}

//
func (this AssertionStatement) Option(option string) string {
	return this.options[option]
}

//
func (this AssertionStatement) Source() Code {
	return ""
}
