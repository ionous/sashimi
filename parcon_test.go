package sashimi

import (
	"fmt"
	"github.com/ionous/sashimi/console"
	"github.com/ionous/sashimi/parser"
	"testing"
)

//
type TestNoun struct {
	name string
}

//
func (this TestNoun) Name() string {
	return this.name
}

// ICommand
type TestNouns struct {
	t     *testing.T
	name  string
	count int
}

func (this TestNouns) NewMatcher() parser.IMatch {
	return &TestMatcher{this.count}
}

func (this TestNouns) RunCommand(nouns ...string) (err error) {
	this.t.Log(">", this.name, nouns)
	if len(nouns) != this.count {
		err = fmt.Errorf("mismatched nouns")
	}
	return err
}

type TestMatcher struct {
	count int
}

func (this *TestMatcher) MatchNoun(word string, article string) (noun string, err error) {
	if this.count <= 0 {
		err = fmt.Errorf("too many nouns")
	} else {
		noun = word
		this.count--
	}
	return noun, err
}

//
func TestConsoleParser(t *testing.T) {

	Look := "look|l"
	Examine := "examine|x|watch|describe|check {{something}}"
	LookThing := `look|l {{something}}`
	LookAt := `look|l at {{something}}`
	Show := "show|present|display {{something}} {{something else}}"
	ShowTo := "show|present|display {{something else}} to {{something}}"

	// commands
	p := parser.NewParser()
	looking, _ := p.AddCommand("looking", TestNouns{t, "l", 0}, 0)
	examining, _ := p.AddCommand("examining it", TestNouns{t, "x", 1}, 1)
	showingItTo, _ := p.AddCommand("showing to it", TestNouns{t, "s", 2}, 2)

	// grammar
	looking.LearnPattern(Look)
	examining.LearnPattern(Examine)
	examining.LearnPattern(LookThing)
	examining.LearnPattern(LookAt)
	showingItTo.LearnPattern(Show)
	showingItTo.LearnPattern(ShowTo)

	// globals
	strs := []string{
		"look",
		"l",
		"look at noun",
		"show robot glory",
		"x one",
		"display success to world",
	}
	c := console.NewBufCon(strs)

	// try it out
	for {
		if s, ok := c.Readln(); !ok {
			break
		} else {
			if name, res, e := p.Parse(s); e != nil {
				t.Fatal(name, e)
			} else {
				res.Run()
			}
		}
	}
}
