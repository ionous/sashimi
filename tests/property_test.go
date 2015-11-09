package tests

import (
	"fmt"
	C "github.com/ionous/sashimi/compiler"
	M "github.com/ionous/sashimi/model"
	. "github.com/ionous/sashimi/script"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"testing"
)

//
// test whether the instance has the value or not
//
func testField(
	res C.MemoryResult,
	instName string,
	fieldName string,
	value interface{},
) (err error) {
	if inst, ok := res.Model.Instances.FindInstance(instName); !ok {
		err = M.InstanceNotFound(instName)
	} else if prop, ok := inst.Class.FindProperty(fieldName); !ok {
		err = fmt.Errorf("'%s.%v' missing field", instName, fieldName)
	} else if val, ok := inst.Value(prop.GetId()); !ok {
		err = fmt.Errorf("'%s.%v' missing value", instName, fieldName)
	} else {
		if enum, ok := prop.(M.EnumProperty); ok {
			val, _ = enum.IndexToChoice(val.(int))
		}
		if !assert.ObjectsAreEqualValues(val, value) {
			err = fmt.Errorf("'%s.%v' %#v!= %#v", instName, fieldName, val, value)
		}
	}
	return err
}

//
// compile nothing succesfully
func TestEmpty(t *testing.T) {
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
func TestClass(t *testing.T) {
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
	stories := res.Model.Classes[M.MakeStringId("stories")]
	if stories == nil {
		t.Fatal("expected stories", res.Model.Classes)
	}
	if stories.Singular != "story" {
		t.Fatal("singular/plural problem", stories)
	}
}

//
// double specification of the same class causes no error
func TestDoubledClass(t *testing.T) {
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
func TestClassProperty(t *testing.T) {
	s := Script{}

	s.The("kinds",
		Have("text", "text"),
		Have("num", "num"),
	)
	if res, err := s.Compile(Log(t)); assert.NoError(t, err) {
		res.Model.PrintModel(t.Log)
		if cls := res.Model.Classes[M.MakeStringId("kinds")]; assert.NotNil(t, cls) {
			if props := cls.Properties; assert.NotNil(t, props) {
				require.Len(t, props, 2+1) // MOD: +1 for auto-generated "name" property
				if pid := M.MakeStringId("text"); assert.Contains(t, props, pid) {
					require.IsType(t, M.TextProperty{}, props[pid])
				}
				if pid := M.MakeStringId("num"); assert.Contains(t, props, pid) {
					require.IsType(t, M.NumProperty{}, props[pid])
				}
			}
		}
	}
}

//
// create an instance ( with no properties )
func TestInst(t *testing.T) {
	s := Script{}
	s.The("kind", Called("test"), Exists())
	res, err := s.Compile(Log(t))
	if err != nil {
		t.Error(err)
	}
	res.Model.PrintModel(t.Log)
	if len(res.Model.Instances) != 1 {
		t.Fatal("expected one instance", res.Model.Instances)
	}
	test := res.Model.Instances[M.MakeStringId("test")]
	if test == nil {
		t.Fatal("expected test instance", res.Model.Instances)
	}
}

//
func TestTextProperties(t *testing.T) {
	s := Script{}

	s.The("kinds",
		Called("stories").WithSingularName("story"),
		Have("author", "text"),
	)
	s.The("story",
		Called("test"),
		Has("author", "any mouse"),
	)
	res, err := s.Compile(Log(t))
	if err != nil {
		t.Fatal(err)
	}
	res.Model.PrintModel(t.Log)
	err = testField(res, "test", "author", "any mouse")
	if err != nil {
		t.Fatal(err)
	}
}

//
func TestNumProperties(t *testing.T) {
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
	err = testField(res, "test", "int", 5)
	if err != nil {
		t.Fatal(err)
	}
	err = testField(res, "test", "float", 3.25)
	if err != nil {
		t.Fatal(err)
	}
}

//
// build several stories with different settings for their yes/no values
// verify all is as expected
func TestEitherOr(t *testing.T) {
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
	res, err := s.Compile(Log(t))
	if err != nil {
		t.Log(err)
	}
	//res.Model.PrintModel(t.Log)
	//
	err = testField(res, "scored-default", "scoredProperty", M.MakeStringId("scored")) // not default: false
	if err != nil {
		t.Fatal(err)
	}
	err = testField(res, "unscored-default", "scoredProperty", M.MakeStringId("unscored")) // not default: false
	if err != nil {
		t.Fatal(err)
	}
	err = testField(res, "scored", "scoredProperty", M.MakeStringId("scored")) // not default: true
	if err != nil {
		t.Fatal(err)
	}
	err = testField(res, "unscored", "scoredProperty", M.MakeStringId("unscored")) // not default: true
	if err != nil {
		t.Fatal(err)
	}
}

//
// choose an unselected value and make sure it reports an error
func TestEitherError(t *testing.T) {
	s := Script{}
	s.The("kinds",
		Called("stories").WithSingularName("story"),
		AreEither("scored").Or("unscored").Usually("unscored"),
	)

	s.The("story",
		Called("errors"),
		Is("this is meant to report an issue"),
	)

	res, err := s.Compile(Log(t))
	if err == nil {
		res.Model.PrintModel(t.Log)
		t.Fatal("expected unscored story to report an issue")
	}
	t.Log("expected error:", err)
}
