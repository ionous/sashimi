package parser

import (
	"fmt"
	"regexp"
	"testing"
)

type NounPool map[string]bool
type TestNounFactory struct {
	t    *testing.T
	pool NounPool
}

// IMatch
func (this *TestNounFactory) MatchNoun(name string, article string) (noun string, err error) {
	if article != "" {
		this.t.Logf("shorted `%s` with article:`%s`", name, article)
	}
	if this.pool[name] {
		noun = name
	} else {
		err = fmt.Errorf("unknown noun %s,%s", article, name)
	}
	return noun, err
}

//
func (this NounPool) AddNouns(names ...string) {
	for _, n := range names {
		this[n] = true
	}
}

// //
func failOnErrors(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
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
	t       *testing.T
	name    string
	expect  []string
	cmd     ILearn
	factory *TestNounFactory
}

func (this *TestCmd) NewMatcher() IMatch {
	return this.factory
}

func (this *TestCmd) RunCommand(nouns ...string) (err error) {
	got, want := len(nouns), len(this.expect)
	if got != want {
		err = fmt.Errorf("noun count doesnt match %d(got)!=%d(want)", got, want)
	} else {
		for i, noun := range nouns {
			want, got := this.expect[i], noun
			if want != got {
				err = fmt.Errorf("nouns dont match %s(got)!=%s(want)", got, want)
				break
			}
		}
	}
	this.t.Log("running", this.name, nouns, err)
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
	factory := &TestNounFactory{t, nil}
	testCmds := []*TestCmd{
		{name: "looking"},
		{name: "examining", expect: []string{"n1"}},
		{name: "showing", expect: []string{"actor", "prize"}},
		{name: "smelling"},
		{name: "spacing", expect: []string{"evil fish"}},
	}

	// ok: add commands
	for _, cmd := range testCmds {
		cmd.t = t
		cmd.factory = factory
		comp, err := p.AddCommand(cmd.name, cmd, len(cmd.expect))
		if err != nil {
			t.Fatal(err, cmd)
		}
		cmd.cmd = comp
	}

	// err: change function
	repeat := testCmds[0]
	_, fail := p.AddCommand(repeat.name, repeat, len(repeat.expect))
	if fail == nil {
		t.Fatalf("expected changed function to fail")
	}

	// ok: learn ya some learnings
	l, x, show, smell, space := testCmds[0], testCmds[1], testCmds[2], testCmds[3], testCmds[4]

	if e := l.cmd.LearnPattern("look|l"); e != nil {
		t.Error(e)
	}
	if e := x.cmd.LearnPattern("examine|x|watch|describe|check {{something}}"); e != nil {
		t.Error(e)
	}
	if e := x.cmd.LearnPattern("look|l {{something}}"); e != nil {
		t.Error(e)
	}
	if e := x.cmd.LearnPattern("look|l at {{something}}"); e != nil {
		t.Error(e)
	}
	if e := show.cmd.LearnPattern("show|present|display {{something}} {{something else}}"); e != nil {
		t.Error(e)
	}
	if e := show.cmd.LearnPattern("show|present|display {{something else}} to {{something}}"); e != nil {
		t.Error(e)
	}
	if e := smell.cmd.LearnPattern("smell"); e != nil {
		t.Error(e)
	}
	if e := space.cmd.LearnPattern("space test {{something}}"); e != nil {
		t.Error(e)
	}

	const (
		ExpectFailure = iota
		ExpectNouns
		ExpectSuccess
	)

	// ok: parse some commands
	// thinking a function for handling them would be good, damn error codes
	testParser := func(cmd string, expect string, expectRes int) {
		normalizedInput := NormalizeInput(cmd)
		found, res, err := p.Parse(normalizedInput)
		expectCmd, foundCmd := expect != "", found != ""
		if expectCmd != foundCmd {
			t.Error(cmd, ": should have been", expectCmd, res, err)
		} else if found != expect {
			t.Error(cmd, ": should have matched", expect, "was", res, err, "instead.")
		} else if foundCmd {
			expectNouns, matchedCmd := expectRes != ExpectFailure, err == nil
			if expectNouns != matchedCmd {
				t.Error(cmd, ": success should have been", expectNouns, res, err)
			} else {
				if err == nil {
					err = res.Run()
				}
				success, expectSuccess := err == nil, expectRes == ExpectSuccess
				if success != expectSuccess {
					t.Error(cmd, ": nouns should have been", ExpectSuccess, res, err)
				}
			}
		}
	}

	nf := make(NounPool)
	factory.pool = nf

	testParser("ignore", "", ExpectFailure)
	testParser("look", "looking", ExpectSuccess)
	testParser("smell", "smelling", ExpectSuccess)
	// should fail because we dont know n1 yet.
	testParser("x n1", "examining", ExpectFailure)
	nf.AddNouns("n1")
	// should succeed because now we know n1.
	testParser("x n1", "examining", ExpectSuccess)
	// should succeed because now we know n1.
	testParser("x the n1", "examining", ExpectSuccess)
	// spaces shouldnt matter
	testParser("  x  n1    ", "examining", ExpectSuccess)
	// we should still fail on some other noun
	testParser("x n2", "examining", ExpectFailure)
	// look something is the same as examining
	testParser("look n1", "examining", ExpectSuccess)
	// look at something is the same as examining
	// unknown noun at n1
	testParser("look at n1", "examining", ExpectSuccess)
	// spaces shouldnt matter
	testParser("look	at	n1", "examining", ExpectSuccess)
	nf.AddNouns("actor", "prize")
	testParser("show actor prize", "showing", ExpectSuccess)
	testParser("present prize to actor", "showing", ExpectSuccess)
	//fail because the test string expects actor first.
	testParser("show prize actor", "showing", ExpectNouns)
	nf.AddNouns("evil fish")
	testParser("space test evil fish", "spacing", ExpectSuccess)
	testParser("space   test   an    evil   fish  ", "spacing", ExpectSuccess)
	// give us some nouns
	testParser("show the actor some prize", "showing", ExpectSuccess)
}
