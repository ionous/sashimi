package parser

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

type NounPool struct {
	pool map[string]bool
}

func (nouns NounPool) FindNoun(name string, article string) (noun string, err error) {
	if nouns.pool[name] {
		noun = name
	} else {
		err = fmt.Errorf("Unknown noun: %s %s.", article, name)
	}
	return noun, err
}

//
func (nouns *NounPool) AddNouns(names ...string) {
	for _, n := range names {
		nouns.pool[n] = true
	}
}

//
func testTokens(t *testing.T, s string, expectGroups, expectTags int) {
	groups, tags := tokenize(s)
	numGroups, numTags := len(groups), len(tags)
	if numGroups != expectGroups || numTags != expectTags {
		t.Errorf("expected %d,%d tokens in `%s` got %d,%d instead.",
			expectGroups, expectTags, s, numGroups, numTags)
		t.Logf("%v %v", groups, tags)
	}
}

//
// make sure we can learn all of the requested tokens
func TestExp(t *testing.T) {
	testTokens(t, "look|l", 1, 0)
	testTokens(t, "{{something}}", 1, 1)
	testTokens(t, "examine|x|watch|describe|check {{something}}", 2, 1)
	testTokens(t, "look|l {{something}}", 2, 1)
	testTokens(t, "look|l at {{something}}", 3, 1)
	testTokens(t, "show|present|display {{something}} {{something else}}", 3, 2)
	testTokens(t, "show|present|display {{something else}} to {{something}}", 4, 2)
}

type TestCmd struct {
	*testing.T
	name   string
	expect []string
	comp   *Comprehension
	pool   NounPool
}

func (cmd *TestCmd) NewMatcher() (IMatch, error) {
	return &TestMatcher{TestCmd: cmd}, nil
}

type TestMatcher struct {
	*TestCmd
	nouns []string
}

func (m *TestMatcher) MatchNoun(word string, article string) (err error) {
	if len(m.nouns) < len(m.expect) {
		if n, e := m.pool.FindNoun(word, article); e != nil {
			err = e
		} else {
			m.nouns = append(m.nouns, n)
		}
	} else {
		err = fmt.Errorf("too many nouns")
	}
	return err
}

// Matched gets called after all nouns in an input have been parsed succesfully.
func (m *TestMatcher) OnMatch() (err error) {
	m.Logf("matched %s, expecting %s", m.nouns, m.expect)
	if got, want := len(m.nouns), len(m.expect); got != want {
		err = fmt.Errorf("noun count doesnt match %d(got)!=%d(want)", got, want)
	} else {
		for i := 0; i < got; i++ {
			got, want := m.nouns[i], m.expect[i]
			if want != got {
				err = fmt.Errorf("nouns dont match %s(got)!=%s(want)", got, want)
				break
			}
		}
	}
	return err
}

func TestLookAt(t *testing.T) {
	look := regexp.MustCompile(`^\s*(?:look|l)\s+(.*?)\s*$`)
	lookAt := regexp.MustCompile(`^\s*(?:look|l)\s+at\s+(.*?)\s*$`)
	lookString := "look n1"
	lookAtString := "look at n1"
	if !look.MatchString(lookString) {
		t.Fatal(lookString)
	}
	if !lookAt.MatchString(lookAtString) {
		t.Fatal(lookAtString)
	}
}

//
func TestUnderstandings(t *testing.T) {

	// create a parser
	p := NewParser()

	// ok: define some commands with their allowed nouns
	pool := NounPool{make(map[string]bool)}
	testCmds := []*TestCmd{
		{name: "looking"},
		{name: "examining", expect: []string{"n1"}},
		{name: "showing", expect: []string{"actor", "prize"}},
		{name: "smelling"},
		{name: "spacing", expect: []string{"evil fish"}},
	}

	// ok: add commands
	for _, cmd := range testCmds {
		cmd.T = t
		cmd.pool = pool
		comp, err := p.NewComprehension(cmd.name, cmd.NewMatcher)
		if err != nil {
			t.Fatal(err, cmd)
		}
		cmd.comp = comp
	}

	// err: change function
	repeat := testCmds[0]
	_, fail := p.NewComprehension(repeat.name, repeat.NewMatcher)
	if fail == nil {
		t.Fatalf("expected changed function to fail")
	}

	// ok: learn ya some learnings
	l, x, show, smell, space := testCmds[0], testCmds[1], testCmds[2], testCmds[3], testCmds[4]

	if _, e := l.comp.LearnPattern("look|l"); e != nil {
		t.Error(e)
	}
	if _, e := x.comp.LearnPattern("examine|x|watch|describe|check {{something}}"); e != nil {
		t.Error(e)
	}
	if _, e := x.comp.LearnPattern("look|l {{something}}"); e != nil {
		t.Error(e)
	}
	if _, e := x.comp.LearnPattern("look|l at {{something}}"); e != nil {
		t.Error(e)
	}
	if _, e := show.comp.LearnPattern("show|present|display {{something}} {{something else}}"); e != nil {
		t.Error(e)
	}
	if _, e := show.comp.LearnPattern("show|present|display {{something else}} to {{something}}"); e != nil {
		t.Error(e)
	}
	if _, e := smell.comp.LearnPattern("smell"); e != nil {
		t.Error(e)
	}
	if _, e := space.comp.LearnPattern("space test {{something}}"); e != nil {
		t.Error(e)
	}

	// ok: parse some commands
	testParser := func(cmd string, expect string) (err error) {
		normalizedInput := NormalizeInput(cmd)
		match, e := p.ParseInput(normalizedInput)
		if e != nil {
			err = e
		} else if match.Pattern != nil && match.Pattern.Comprehension().Name() != expect {
			err = fmt.Errorf("Mismatched pattern: %s got %s", expect, match.Pattern)
		} else {
			err = match.OnMatch()
		}
		return err
	}

	nf := pool

	// assert.Error(t, testParser("ignore", ""), "doesnt exist")
	// assert.NoError(t, testParser("look", "looking"), "")
	// assert.NoError(t, testParser("smell", "smelling"), "")
	// assert.Error(t, testParser("x n1", "examining"), "we dont know n1 yet.")
	nf.AddNouns("n1")
	// assert.NoError(t, testParser("x n1", "examining"), "We should know n1.")
	// assert.NoError(t, testParser("x the n1", "examining"), "now we know n1 with an article")
	// assert.NoError(t, testParser("  x  n1    ", "examining"), "spaces shouldnt matter")
	// assert.Error(t, testParser("x n2", "examining"), "fail on another unknown noun")
	// assert.NoError(t, testParser("look n1", "examining"), "looking is examining")
	// assert.NoError(t, testParser("look at n1", "examining"), "look at is examining")
	// assert.NoError(t, testParser("look	at	n1", "examining"), "ignore spaces")
	nf.AddNouns("actor", "prize")
	// assert.NoError(t, testParser("show actor prize", "showing"), "test showing")
	// assert.NoError(t, testParser("present prize to actor", "showing"), "reverse showing")
	assert.Error(t, testParser("show prize actor", "showing"), "because the test string expects actor first.")
	nf.AddNouns("evil fish")
	// assert.NoError(t, testParser("space test evil fish", "spacing"), "spacing in nouns")
	// assert.NoError(t, testParser("space   test   an    evil   fish  ", "spacing"), "fishy spacing")
	// assert.NoError(t, testParser("show the actor some prize", "showing"), "give us some nouns")
}
