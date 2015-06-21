package script

import (
	"regexp"
	"strings"
)

//
// Fragment to assert the existence of a class or instance
// this.The("room", this.Called("home" )
// this.The("kinds", this.Called("coins").Singular("coin")
func Called(subject string) CalledFragment {
	return CalledFragment{NewOrigin(1), subject, ""}
}

//
func (this CalledFragment) WithSingularName(name string) IFragment {
	this.singular = name
	return this
}

//
// Exists(): optional fragment for making otherwise empty statements more readable
// The("room", Called("parlor of despair"), Exists())
//
func Exists() IFragment {
	return FunctionFragment{
		func(b SubjectBlock) error {
			return nil
		}}
}
func Exist() IFragment { return Exists() }

//
// implementation:
//
type CalledFragment struct {
	origin   Origin
	subject  string // name of the class or instance being declared
	singular string // optional singular version of that name
}

var articles = regexp.MustCompile(`^((?U)the|a|an|our|some) `)

func (this CalledFragment) MakeStatement(b SubjectBlock) error {
	// FIX: this is only half measure --
	// really it needs to split into words, then compare the first article.
	name := strings.TrimSpace(this.subject)
	article, bare := "", name
	pair := articles.FindStringIndex(name)
	if pair != nil {
		article = name[:pair[0]]
		bare = name[pair[1]:]
	}
	opt := map[string]string{
		"article":       article,
		"long name":     name,
		"singular name": this.singular,
	}
	return b.NewAssertion(b.theKeyword, bare, opt)
}
