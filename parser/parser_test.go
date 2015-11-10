package parser

import (
	"fmt"
	"github.com/ionous/sashimi/util/ident"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

type NounPool struct {
	nouns map[string]bool
}

func (np NounPool) FindNoun(name string, article string) (noun string, err error) {
	if np.nouns[name] {
		noun = name
	} else {
		err = fmt.Errorf("Unknown noun: %s %s.", article, name)
	}
	return noun, err
}

//
func (np *NounPool) AddNouns(names ...string) {
	for _, n := range names {
		np.nouns[n] = true
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
	nouns  NounPool
}

type CommandPool map[ident.Id]*TestCmd

func (cmds CommandPool) NewMatcher(id ident.Id) (ret IMatch, err error) {
	if cmd, ok := cmds[id]; !ok {
		err = fmt.Errorf("cmd %s not found", id)
	} else {
		ret = &TestMatcher{cmd: cmd}
	}
	return
}

type TestMatcher struct {
	cmd   *TestCmd
	nouns []string
}

func (m *TestMatcher) MatchNoun(word string, article string) (err error) {
	if len(m.nouns) < len(m.cmd.expect) {
		if n, e := m.cmd.nouns.FindNoun(word, article); e != nil {
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
	m.cmd.Logf("matched %s, expecting %s", m.nouns, m.cmd.expect)
	if got, want := len(m.nouns), len(m.cmd.expect); got != want {
		err = fmt.Errorf("noun count doesnt match %d(got)!=%d(want)", got, want)
	} else {
		for i := 0; i < got; i++ {
			got, want := m.nouns[i], m.cmd.expect[i]
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
	commandPool := make(CommandPool)
	p := NewParser(commandPool)

	// ok: define some commands with their allowed nouns
	nouns := NounPool{make(map[string]bool)}

	testCmds := []*TestCmd{
		{name: "looking"},
		{name: "examining", expect: []string{"n1"}},
		{name: "reporting shown", expect: []string{"actor", "prize"}},
		{name: "smelling"},
		{name: "spacing", expect: []string{"evil fish"}},
	}

	// ok: add commands
	for _, cmd := range testCmds {
		id := ident.MakeId(cmd.name)
		cmd.T = t
		cmd.nouns = nouns
		comp, err := p.NewComprehension(id)
		if err != nil {
			t.Fatal(err, cmd)
		}
		cmd.comp = comp
		commandPool[id] = cmd
	}

	// err: change function
	repeat := testCmds[0]
	_, fail := p.NewComprehension(ident.MakeId(repeat.name))
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
		if pattern, matcher, e := p.ParseInput(normalizedInput); e != nil {
			err = e
		} else if pattern != nil && pattern.Comprehension().Id() != ident.MakeId(expect) {
			err = fmt.Errorf("Mismatched pattern: %s got %s", expect, pattern)
		} else if m, ok := matcher.(*TestMatcher); assert.True(t, ok) {
			err = m.OnMatch()
		}
		return err
	}

	nf := nouns

	assert.Error(t, testParser("ignore", ""), "doesnt exist")
	assert.NoError(t, testParser("look", "looking"), "")
	assert.NoError(t, testParser("smell", "smelling"), "")
	assert.Error(t, testParser("x n1", "examining"), "we dont know n1 yet.")
	nf.AddNouns("n1")
	assert.NoError(t, testParser("x n1", "examining"), "We should know n1.")
	assert.NoError(t, testParser("x the n1", "examining"), "now we know n1 with an article")
	assert.NoError(t, testParser("  x  n1    ", "examining"), "spaces shouldnt matter")
	assert.Error(t, testParser("x n2", "examining"), "fail on another unknown noun")
	assert.NoError(t, testParser("look n1", "examining"), "looking is examining")
	assert.NoError(t, testParser("look at n1", "examining"), "look at is examining")
	assert.NoError(t, testParser("look	at	n1", "examining"), "ignore spaces")
	nf.AddNouns("actor", "prize")
	assert.NoError(t, testParser("show actor prize", "reporting shown"), "test showing")
	assert.NoError(t, testParser("present prize to actor", "reporting shown"), "reverse showing")
	assert.Error(t, testParser("show prize actor", "reporting shown"), "because the test string expects actor first.")
	nf.AddNouns("evil fish")
	assert.NoError(t, testParser("space test evil fish", "spacing"), "spacing in nouns")
	assert.NoError(t, testParser("space   test   an    evil   fish  ", "spacing"), "fishy spacing")
	assert.NoError(t, testParser("show the actor some prize", "reporting shown"), "give us some nouns")
}
