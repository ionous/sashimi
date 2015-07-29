package script

import (
	S "github.com/ionous/sashimi/source"
	"regexp"
	"strings"
)

//
// Fragment to assert the existence of a class or instance
// The("room", Called("home"))
func Called(subject string) CalledFragment {
	origin := NewOrigin(2)
	return CalledFragment{origin, subject, ""}
}

//
func (frag CalledFragment) WithSingularName(name string) IFragment {
	frag.singular = name
	return frag
}

//
// Exists(): optional fragment for making otherwise empty statements more readable
// The("room", Called("parlor of despair"), Exists())
//
func Exists() IFragment {
	return NewFunctionFragment(func(b SubjectBlock) error {
		return nil
	})
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

func (frag CalledFragment) MakeStatement(b SubjectBlock) error {
	// FIX: this is only half measure --
	// really it needs to split into words, then compare the first article.
	name := strings.TrimSpace(frag.subject)
	article, bare := "", name
	if pair := articles.FindStringIndex(name); pair != nil {
		article = name[:pair[0]]
		bare = name[pair[1]:]
	}
	opt := map[string]string{
		"article":       article,
		"long name":     name,
		"singular name": frag.singular,
	}
	fields := S.AssertionFields{b.theKeyword, bare, opt}
	return b.NewAssertion(fields, frag.origin.Code())
}
