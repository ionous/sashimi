package tests

import (
	. "github.com/ionous/mars/script"
	"github.com/ionous/sashimi/compiler"
	M "github.com/ionous/sashimi/compiler/model"
	S "github.com/ionous/sashimi/source"
	"github.com/ionous/sashimi/util/errutil"
	"github.com/ionous/sashimi/util/ident"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
)

// compile nothing succesfully
func TestPropertyEmpty(t *testing.T) {
	s := NewScript()
	src := &S.Statements{}
	if e := s.Generate(src); assert.NoError(t, e) {
		if model, e := compiler.Compile(*src, ioutil.Discard); assert.NoError(t, e) {
			// we expect the single built in "kinds" class
			if assert.Len(t, model.Classes, 1, "expected one classes") {
				return
			}
		}
	}
	t.FailNow()
}

// create a single subclass called stories
func TestPropertySubclass(t *testing.T) {
	s := The("kinds",
		Called("stories"), HasSingularName("story"),
	)
	src := &S.Statements{}
	if e := s.Generate(src); assert.NoError(t, e) {
		if model, e := compiler.Compile(*src, ioutil.Discard); assert.NoError(t, e) {
			if assert.Len(t, model.Classes, 2, "expected two classes") {
				if stories := model.Classes[ident.MakeId("stories")]; assert.NotNil(t, stories, "expected stories") {
					if assert.EqualValues(t, stories.Singular, "story", "singular/plural problem") {
						return
					}
				}
			}
		}
	}
	t.FailNow()
}

// double specification of the same class causes no error
func TestPropertyDoubledClass(t *testing.T) {
	s := NewScript(
		The("kinds", Called("stories"), HasSingularName("story")),
		The("kinds", Called("stories"), Exists()),
	)
	src := &S.Statements{}
	if e := s.Generate(src); assert.NoError(t, e) {
		if _, e := compiler.Compile(*src, ioutil.Discard); assert.NoError(t, e) {
			return
		}
	}
	t.FailNow()
}

// create properties on the built in object class
func TestPropertyKinds(t *testing.T) {
	s := The("kinds",
		Have("text", "text"),
		Have("num", "num"),
	)
	src := &S.Statements{}
	if e := s.Generate(src); assert.NoError(t, e) {
		if model, e := compiler.Compile(*src, ioutil.Discard); assert.NoError(t, e) {
			if cls := model.Classes[ident.MakeId("kinds")]; assert.NotNil(t, cls) {
				if props := cls.Properties; assert.NotNil(t, props) {
					require.Len(t, props, 2+1) // MOD: +1 for auto-generated "name" property
					if p, ok := cls.FindProperty("text"); assert.True(t, ok, "found text") {
						require.EqualValues(t, M.TextProperty, p.Type)
					}
					if p, ok := cls.FindProperty("num"); assert.True(t, ok, "found num") {
						require.EqualValues(t, M.NumProperty, p.Type)
					}
					return
				}
			}
		}
	}
	t.FailNow()
}

// TestPropertyInst: create an instance ( with no properties )
func TestPropertyInst(t *testing.T) {
	s := The("kind", Called("test"), Exists())
	src := &S.Statements{}
	if e := s.Generate(src); assert.NoError(t, e) {
		if model, err := compiler.Compile(*src, ioutil.Discard); assert.NoError(t, err, "compile") {
			//	model.PrintModel(t.Log)
			if assert.Len(t, model.Instances, 1, "expected one instance") {
				if test, ok := model.Instances[ident.MakeId("test")]; assert.True(t, ok, "expected test instance") {
					// test auto-generated name.
					nameId := ident.Join(ident.MakeId("kinds"), ident.MakeId("name"))
					if name, ok := test.Values[nameId]; assert.True(t, ok, "have name value") {
						require.EqualValues(t, "test", name)
					}
					return
				}
			}
		}
	}
	t.FailNow()
}

func TestPropertyText(t *testing.T) {
	s := NewScript(
		The("kinds",
			Called("stories"), HasSingularName("story"),
			Have("author", "text"),
		),
		The("story",
			Called("test"),
			Has("author", "any mouse"),
		))
	src := &S.Statements{}
	if e := s.Generate(src); assert.NoError(t, e) {
		if model, err := compiler.Compile(*src, ioutil.Discard); assert.NoError(t, err, "compile") {
			if v, err := field(model, "test", "author"); assert.NoError(t, err, "test field") {
				require.EqualValues(t, "any mouse", v, "mismatched")
				return
			}
		}
	}
	t.FailNow()
}

func TestPropertyNum(t *testing.T) {
	s := NewScript(
		The("kinds",
			Called("stories"), HasSingularName("story"),
			Have("int", "num"),
			Have("float", "num")),
		The("story",
			Called("test"),
			Has("int", 5)),
		The("test",
			Has("float", 3.25)),
	)
	src := &S.Statements{}
	if e := s.Generate(src); assert.NoError(t, e) {
		if model, e := compiler.Compile(*src, ioutil.Discard); assert.NoError(t, e) {
			if v, e := field(model, "test", "int"); assert.NoError(t, e) {
				require.EqualValues(t, 5, v, "int mismatch")
				if v, e := field(model, "test", "float"); assert.NoError(t, e) {
					require.EqualValues(t, 3.25, v, "float mismatch")
					return
				}
			}
		}
	}
	t.FailNow()
}

// build several stories with different settings for their yes/no values
// verify all is as expected
// go test -run TestPropertyEitherOr
func TestPropertyEitherOr(t *testing.T) {
	s := NewScript(
		The("kinds",
			Called("scored-stories"), HasSingularName("scored-story"),
			AreEither("scored").Or("unscored"),
		),
		The("kinds",
			Called("unscored-stories"), HasSingularName("unscored-story"),
			AreOneOf("scored", "unscored").Usually("unscored"),
		),
		The("scored-story",
			Called("scored-default"),
			//Is("scored"),
		),
		The("unscored-story",
			Called("unscored-default"),
			//Is("unscored"),
		),
		The("unscored-story",
			Called("scored"),
			Is("scored"),
		),
		The("unscored-story",
			Called("unscored"),
			Is("unscored"),
		),
	)
	src := &S.Statements{}
	if e := s.Generate(src); assert.NoError(t, e) {
		if model, e := compiler.Compile(*src, ioutil.Discard); assert.NoError(t, e) {
			if v, e := field(model, "scored-default", "scored"); assert.NoError(t, e) {
				// not default: false
				require.EqualValues(t, "scored", v)
				if v, e := field(model, "unscored-default", "scored"); assert.NoError(t, e) {
					// not default: false
					require.EqualValues(t, "unscored", v)
					if v, e := field(model, "scored", "scored"); assert.NoError(t, e) {
						// not default: true
						require.EqualValues(t, "scored", v)
						if v, e := field(model, "unscored", "scored"); assert.NoError(t, e) {
							// not default: true
							require.EqualValues(t, "unscored", v)
							return
						}
					}
				}
			}
		}
	}
	t.FailNow()
}

// choose an unselected value and make sure it reports an error
func TestPropertyEitherError(t *testing.T) {
	s := NewScript(
		The("kinds",
			Called("stories"), HasSingularName("story"),
			AreEither("scored").Or("unscored").Usually("unscored"),
		),
		The("story",
			Called("errors"),
			Is("this is meant to report an issue"),
		),
	)
	src := &S.Statements{}
	if e := s.Generate(src); assert.NoError(t, e) {
		if _, e := compiler.Compile(*src, ioutil.Discard); assert.Error(t, e, "expects compile failure") {
			return
		}
	}
	t.FailNow()
}

// test whether the instance has the value or not
func field(
	model compiler.MemoryResult,
	instName string,
	fieldName string,
) (ret interface{}, err error) {
	if inst, ok := model.Instances[ident.MakeId(instName)]; !ok {
		err = errutil.New("instance not found", instName)
	} else if cls, ok := model.Classes[inst.Class]; !ok {
		err = errutil.New("class not found for", inst.Class, instName)
	} else if prop, ok := cls.FindProperty(fieldName); !ok {
		err = errutil.New("missing field", instName, fieldName)
	} else {
		val, valExisted := inst.Values[prop.Id]

		if prop.Type == M.EnumProperty {
			if enum, ok := model.Enumerations[prop.Id]; !ok {
				err = errutil.New("missing enum", instName, prop.Id)
			} else if !valExisted {
				val = enum.Best()
				valExisted = true
			}

		}
		if err == nil {
			if !valExisted {
				err = errutil.New("missing value", instName, prop.Id)
			} else {
				ret = val
			}
		}
	}
	return
}
