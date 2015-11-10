package tests

import (
	"fmt"
	"github.com/ionous/sashimi/console"
	"github.com/ionous/sashimi/parser"
	"github.com/ionous/sashimi/util/ident"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

// IComprehend
type TestComprehension struct {
	*testing.T
	name  string
	count int
}

// IMatch
type TestMatcher struct {
	test  TestComprehension
	nouns []string
}

func (m *TestMatcher) MatchNoun(word string, article string) (err error) {
	if len(m.nouns) >= m.test.count {
		err = fmt.Errorf("too many nouns")
	} else {
		m.nouns = append(m.nouns, word)
	}
	return err
}

func (m *TestMatcher) OnMatch() (err error) {
	m.test.Log(">", m.test.name, m.nouns)
	if len(m.nouns) != m.test.count {
		err = fmt.Errorf("mismatched nouns")
	}
	return err
}

type TestPool struct {
	comps parser.Comprehensions
	tests map[ident.Id]TestComprehension
}

func (t TestPool) NewMatcher(id ident.Id) (ret parser.IMatch, err error) {
	if test, ok := t.tests[id]; !ok {
		err = fmt.Errorf("couldnt find test %s", id)
	} else {
		ret = &TestMatcher{test: test}
	}
	return
}

func (t TestPool) NewComprehension(name string, x TestComprehension) (ret *parser.Comprehension, err error) {
	id := ident.MakeId(name)
	if r, e := t.comps.NewComprehension(id); e != nil {
		err = e
	} else {
		ret = r
		t.tests[id] = x
	}
	return
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
	p := TestPool{make(parser.Comprehensions), make(map[ident.Id]TestComprehension)}
	par := parser.P{p, p.comps}

	looking, e := p.NewComprehension("looking", TestComprehension{t, "l", 0})
	require.NoError(t, e)
	examining, e := p.NewComprehension("examining it", TestComprehension{t, "x", 1})
	require.NoError(t, e)
	showingItTo, e := p.NewComprehension("showing to it", TestComprehension{t, "s", 2})
	require.NoError(t, e)
	require.Len(t, p.comps, 3)

	// grammar
	_, e = looking.LearnPattern(Look)
	require.NoError(t, e)
	_, e = examining.LearnPattern(Examine)
	require.NoError(t, e)
	_, e = examining.LearnPattern(LookThing)
	require.NoError(t, e)
	_, e = examining.LearnPattern(LookAt)
	require.NoError(t, e)
	_, e = showingItTo.LearnPattern(Show)
	require.NoError(t, e)
	_, e = showingItTo.LearnPattern(ShowTo)
	require.NoError(t, e)

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
		} else if _, m, err := par.ParseInput(s); assert.NoError(t, err, s) {
			if err := m.(*TestMatcher).OnMatch(); assert.NoError(t, err, s) {
				continue
			}
		}
	}
}
