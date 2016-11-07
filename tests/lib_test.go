package tests

import (
	"github.com/ionous/mars"
	"github.com/ionous/mars/rt"
	. "github.com/ionous/mars/script"
	"github.com/ionous/mars/std"
	S "github.com/ionous/sashimi/source"
	"github.com/ionous/sashimi/util/errutil"
	"github.com/stretchr/testify/require"
	"testing"
)

// Arc implements the script Trytime interface
type Arc struct {
	test *TestGame
}

func (a *Arc) Parse(in string) (ret string, err error) {
	if outs, e := a.test.RunInput(in); e != nil {
		err = errutil.New("error running input:", e)
	} else {
		ret = outs[0]
	}
	return
}

func (a *Arc) Execute(ex rt.Execute) (ret string, err error) {
	if e := ex.Execute(a.test.Game); e != nil {
		err = errutil.New("error testing execution", e)
	} else if out, e := a.test.FlushOutput(); e != nil {
		err = errutil.New("error flushing output", e)
	} else {
		ret = out[0]
	}
	return
}

func (a *Arc) Test(be rt.BoolEval) (err error) {
	if b, e := be.GetBool(a.test.Game); e != nil {
		err = errutil.New("error testing boolean", e)
	} else if !b {
		err = errutil.New("test failed")
	}
	return
}

func libTest(t *testing.T, lib *mars.Package, base *S.Statements, parser string) (err error) {
	if e := lib.Generate(base); e != nil {
		err = errutil.New("error generating lib source", e)
	} else {
		// FIX? serialize the test scripts?
		for _, suite := range lib.Tests {
			src := *base
			if e := suite.Setup.Generate(&src); e != nil {
				err = errutil.New("error generating test suite:", e)
				break
			} else if test, e := NewTestGameSource(t, src, parser, nil); e != nil {
				err = errutil.New("error creating game:", e)
				break
			} else if e := suite.Test(&Arc{&test}); e != nil {
				err = errutil.New("error testing lib:", e)
				break
			}
		}
	}
	return err
}

// func TestLibLang(t *testing.T) {
// base := &S.Statements{}
// 	The("kind", Called("no parser")).Generate(base)
// 	require.NoError(t, libTest(t, &lang.Lang, base, "no parser"))
// }

func TestLibStd(t *testing.T) {
	base := &S.Statements{}
	script := The("actor", Called("player"), Exists())
	script.Generate(base)
	require.NoError(t, libTest(t, &std.Std, base, "player"))
}
