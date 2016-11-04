package tests

import (
	//	"github.com/ionous/mars"
	"github.com/ionous/mars/lang"
	"github.com/ionous/mars/rt"
	. "github.com/ionous/mars/script"
	S "github.com/ionous/sashimi/source"
	"github.com/ionous/sashimi/util/errutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type Arc struct {
	test *TestGame
}

func (a *Arc) Parse(in string) (ret string, err error) {
	if outs, e := a.test.RunInput(in); e != nil {
		err = e
	} else {
		ret = outs[0]
	}
	return
}

func (a *Arc) Execute(ex rt.Execute) (ret string, err error) {
	if e := ex.Execute(a.test.Game.Rtm); e != nil {
		err = e
	} else if out, e := a.test.FlushOutput(); e != nil {
		err = e
	} else {
		ret = out[0]
	}
	return
}

func (a *Arc) Test(be rt.BoolEval) (err error) {
	if b, e := be.GetBool(a.test.Game.Rtm); e != nil {
		err = e
	} else if !b {
		err = errutil.New("test failed")
	}
	return
}

func TestLibLang(t *testing.T) {
	base := &S.Statements{}
	The("kind", Called("no parser")).Generate(base)

	lib := &lang.Lang
	if err := lib.Generate(base); assert.NoError(t, err, "build") {
		// FIX? serialize the test scripts?
		for _, suite := range lib.Tests {
			// s.Setup -> contains a bunch of specs we have to compile.
			t.Log("testing suite", suite)
			src := *base
			suite.Setup.Generate(&src)

			// FIX FIX FIX --
			// have to add lang scripts and all of its imports -- and can exclude dependenices
			// an interface which wraps -- possibly tests wraps the main packaeg
			if test, err := NewTestGameSource(t, src, "no parser", nil); assert.NoError(t, err) {
				arc := &Arc{&test}

				for _, trial := range suite.Trials {
					err := trial.Test(arc)
					require.NoError(t, err, trial.Name)
				}
			}
		}
	}
}
