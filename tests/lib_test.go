package tests

import (
	"flag"
	"github.com/ionous/mars"
	"github.com/ionous/mars/core"
	"github.com/ionous/mars/lang"
	"github.com/ionous/mars/rt"
	. "github.com/ionous/mars/script"
	"github.com/ionous/mars/std"
	"github.com/ionous/sashimi/meta"
	S "github.com/ionous/sashimi/source"
	"github.com/ionous/sashimi/util/errutil"
	"github.com/ionous/sashimi/util/sbuf"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Arc implements the script Trytime interface
type Arc struct {
	test *TestGame
}

func (a *Arc) Parse(in string) (ret []string, err error) {
	if out, e := a.test.RunInput(in); e != nil {
		err = errutil.New("error running input:", e)
	} else {
		ret = out
	}
	return
}

func (a *Arc) Run(in string, args []meta.Generic) (ret []string, err error) {
	if e := a.test.RunNamedAction(in, args...); e != nil {
		err = errutil.New("run", e)
	} else if out, e := a.test.FlushOutput(); e != nil {
		err = errutil.New("run flush", e)
	} else {
		ret = out
	}
	return
}

func (a *Arc) Execute(ex rt.Execute) (ret []string, err error) {
	if e := ex.Execute(a.test.Game); e != nil {
		err = errutil.New("execute", e)
	} else if out, e := a.test.FlushOutput(); e != nil {
		err = errutil.New("execute flush", e)
	} else {
		ret = out
	}
	return
}

func (a *Arc) Test(be rt.BoolEval) (err error) {
	if b, e := be.GetBool(a.test.Game); e != nil {
		err = errutil.New("testing boolean failed because", e)
	} else if !b.Value {
		err = errutil.New("test failed")
	}
	return
}

func libTest(t *testing.T, lib *mars.Package, name string, base *S.Statements, parser string, pc ParentCreator) (err error) {
	// FIX? serialize the test scripts?
	for _, suite := range lib.Tests {
		if name == "*" || name == suite.Name {
			t.Log(suite.Name)
			for _, unit := range suite.Units {
				src := *base
				if e := lib.Generate(&src); e != nil {
					err = errutil.New("Error generating lib source", sbuf.Q(suite.Name), e)
				} else {
					//pretty.Println(src)
					if e := unit.Setup.Generate(&src); e != nil {
						err = errutil.New("Error generating test suite", sbuf.Q(suite.Name), e)
						break
					} else if test, e := NewTestGameSource(t, src, parser, pc); e != nil {
						err = errutil.New("Error creating game", sbuf.Q(suite.Name), e)
						break
					} else if e := unit.Test(&Arc{&test}); e != nil {
						err = errutil.New("Error testing lib", sbuf.Q(suite.Name), e)
						break
					}
				}
			}
		}
	}
	return err
}

var testName = "*"

// go test -run LibS -v -named Debugging
func init() {
	flag.StringVar(&testName, "named", testName, "select sub test")
	flag.Parse()
}

func TestLibCore(t *testing.T) {
	base := &S.Statements{}
	The("kind", Called("no parser")).Generate(base)
	assert.NoError(t, libTest(t, core.Core(), testName, base, "no parser", nil))
}

func TestLibLang(t *testing.T) {
	base := &S.Statements{}
	The("kind", Called("no parser")).Generate(base)
	assert.NoError(t, libTest(t, lang.Lang(), testName, base, "no parser", nil))
}

func TestLibStd(t *testing.T) {
	base := &S.Statements{}
	script := The("actor", Called("player"), Exists())
	script.Generate(base)
	assert.NoError(t, libTest(t, std.Std(), testName, base, "player", NewStandardParents))
}
