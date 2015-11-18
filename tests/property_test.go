package tests

import (
	"fmt"
	C "github.com/ionous/sashimi/compiler"
	M "github.com/ionous/sashimi/model"
	. "github.com/ionous/sashimi/script"
	"github.com/ionous/sashimi/util/ident"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

//
// test whether the instance has the value or not
//
func testField(
	t *testing.T,
	res C.MemoryResult,
	instName string,
	fieldName string,
	value interface{},
) {
	inst, ok := res.Model.Instances[ident.MakeId(instName)]
	require.True(t, ok, fmt.Sprintf("instance %s not found", instName))

	cls, ok := res.Model.Classes[inst.Class]
	require.True(t, ok, fmt.Sprintf("class %s not found for %s", inst.Class, instName))

	prop, ok := cls.FindProperty(fieldName)
	require.True(t, ok, fmt.Sprintf("'%s.%v' missing field", instName, fieldName))

	val, valExisted := inst.Values[prop.Id]

	if prop.Type == M.EnumProperty {
		if enum, ok := res.Model.Enumerations[prop.Id]; assert.True(t, ok, fmt.Sprintf("'%s.%v' missing enum", instName, prop.Id)) {
			if valExisted {
				val = enum.IndexToChoice(val.(int))
			} else {
				val = enum.Choices[0]
				valExisted = true
			}
		}
	}
	//
	if assert.True(t, valExisted, fmt.Sprintf("'%s.%v' missing value", instName, prop.Id)) {
		require.EqualValues(t, val, value, fmt.Sprintf("'%s.%v'", instName, fieldName))
	}
}

// compile nothing succesfully
func TestPropertyEmpty(t *testing.T) {
	s := Script{}
	c, err := s.Compile(Log(t))
	if err != nil {
		t.Error(err)
	}
	// we expect the single built in "kinds" class
	if len(c.Model.Classes) != 1 {
		t.Fatal("expected one classes", c.Model.Classes)
	}
}

//
// create a single subclass called stories
func TestPropertySubclass(t *testing.T) {
	s := Script{}
	s.The("kinds",
		Called("stories").WithSingularName("story"),
	)
	res, err := s.Compile(Log(t))
	// no expected errors
	if err != nil {
		t.Error(err)
	}
	res.Model.PrintModel(t.Log)
	if len(res.Model.Classes) != 2 {
		t.Fatal("expected two classes", res.Model.Classes)
	}
	stories := res.Model.Classes[ident.MakeId("stories")]
	if stories == nil {
		t.Fatal("expected stories", res.Model.Classes)
	}
	if stories.Singular != "story" {
		t.Fatal("singular/plural problem", stories)
	}
}

//
// double specification of the same class causes no error
func TestPropertyDoubledClass(t *testing.T) {
	s := Script{}
	s.The("kinds",
		Called("stories").WithSingularName("story"),
	)
	s.The("kinds",
		Called("stories"))
	_, err := s.Compile(Log(t))
	if err != nil {
		t.Error(err)
	}
}

//
// create properties on the built in object class
// ( these use extensions )
func TestPropertyKinds(t *testing.T) {
	s := Script{}

	s.The("kinds",
		Have("text", "text"),
		Have("num", "num"),
	)
	if res, err := s.Compile(Log(t)); assert.NoError(t, err) {
		res.Model.PrintModel(t.Log)
		if cls := res.Model.Classes[ident.MakeId("kinds")]; assert.NotNil(t, cls) {
			if props := cls.Properties; assert.NotNil(t, props) {
				require.Len(t, props, 2+1) // MOD: +1 for auto-generated "name" property
				if p, ok := cls.FindProperty("text"); assert.True(t, ok, "found text") {
					require.EqualValues(t, M.TextProperty, p.Type)
				}
				if p, ok := cls.FindProperty("num"); assert.True(t, ok, "found num") {
					require.EqualValues(t, M.NumProperty, p.Type)
				}
			}
		}
	}
}

// TestPropertyInst: create an instance ( with no properties )
// go test -run TestPropertyInst
func TestPropertyInst(t *testing.T) {
	s := Script{}
	s.The("kind", Called("test"), Exists())
	if res, err := s.Compile(Log(t)); assert.NoError(t, err, "compile") {
		//	res.Model.PrintModel(t.Log)
		if assert.Len(t, res.Model.Instances, 1, "expected one instance") {
			if test, ok := res.Model.Instances[ident.MakeId("test")]; assert.True(t, ok, "expected test instance") {
				// test auto-generated name.
				nameId := ident.Join(ident.MakeId("kinds"), ident.MakeId("name"))
				if name, ok := test.Values[nameId]; assert.True(t, ok, "have name value") {
					require.Equal(t, "test", name)
				}
			}
		}
	}
}

//
func TestPropertyText(t *testing.T) {
	s := Script{}

	s.The("kinds",
		Called("stories").WithSingularName("story"),
		Have("author", "text"),
	)
	s.The("story",
		Called("test"),
		Has("author", "any mouse"),
	)
	if res, err := s.Compile(Log(t)); assert.NoError(t, err, "compile") {
		res.Model.PrintModel(t.Log)
		testField(t, res, "test", "author", "any mouse")
	}
}

//
func TestPropertyNum(t *testing.T) {
	s := Script{}

	s.The("kinds",
		Called("stories").WithSingularName("story"),
		Have("int", "num"),
		Have("float", "num"))
	s.The("story",
		Called("test"),
		Has("int", 5))
	s.The("test",
		Has("float", 3.25))
	res, err := s.Compile(Log(t))
	if err != nil {
		t.Fatal(err)
	}
	res.Model.PrintModel(t.Log)
	testField(t, res, "test", "int", 5)
	testField(t, res, "test", "float", 3.25)
}

// build several stories with different settings for their yes/no values
// verify all is as expected
// go test -run TestPropertyEitherOr
func TestPropertyEitherOr(t *testing.T) {
	s := Script{}

	s.The("kinds",
		Called("scored-stories").WithSingularName("scored-story"),
		AreEither("scored").Or("unscored"),
	)
	s.The("kinds",
		Called("unscored-stories").WithSingularName("unscored-story"),
		AreEither("scored").Or("unscored").Usually("unscored"),
	)
	s.The("scored-story",
		Called("scored-default"),
		//Is("scored"),
	)
	s.The("unscored-story",
		Called("unscored-default"),
		//Is("unscored"),
	)
	s.The("unscored-story",
		Called("scored"),
		Is("scored"),
	)
	s.The("unscored-story",
		Called("unscored"),
		Is("unscored"),
	)
	if res, err := s.Compile(Log(t)); assert.NoError(t, err, "compile") {
		res.Model.PrintModel(t.Log)
		//
		testField(t, res, "scored-default", "scored", ident.MakeId("scored"))     // not default: false
		testField(t, res, "unscored-default", "scored", ident.MakeId("unscored")) // not default: false
		testField(t, res, "scored", "scored", ident.MakeId("scored"))             // not default: true
		testField(t, res, "unscored", "scored", ident.MakeId("unscored"))         // not default: true
	}
}

//
// choose an unselected value and make sure it reports an error
func TestPropertyEitherError(t *testing.T) {
	s := Script{}
	s.The("kinds",
		Called("stories").WithSingularName("story"),
		AreEither("scored").Or("unscored").Usually("unscored"),
	)

	s.The("story",
		Called("errors"),
		Is("this is meant to report an issue"),
	)
	if res, err := s.Compile(Log(t)); assert.Error(t, err, "expects compile failure") {
		res.Model.PrintModel(t.Log)
	}
}
