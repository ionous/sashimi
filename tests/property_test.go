package tests

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
	. "github.com/ionous/sashimi/script"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

//
// test whether the instance has the value or not
//
func testField(
	res *M.Model,
	instName string,
	fieldName string,
	value interface{},
) (err error) {
	if inst, ok := res.Instances.FindInstance(instName); !ok {
		err = M.InstanceNotFound(instName)
	} else if prop, ok := inst.Class().FindProperty(fieldName); !ok {
		err = fmt.Errorf("'%s.%v' missing field", instName, fieldName)
	} else if val, ok := inst.Value(prop.Id()); !ok {
		err = fmt.Errorf("'%s.%v' missing value", instName, fieldName)
	} else {
		if enum, ok := prop.(*M.EnumProperty); ok {
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
	c, err := s.Compile(os.Stderr)
	if err != nil {
		t.Error(err)
	}
	// we expect the single built in "kinds" class
	if len(c.Classes) != 1 {
		t.Fatal("expected one classes", c.Classes)
	}
}

//
// create a single subclass called stories
func TestClass(t *testing.T) {
	s := Script{}
	s.The("kinds",
		Called("stories").WithSingularName("story"),
	)
	res, err := s.Compile(os.Stderr)
	// no expected errors
	if err != nil {
		t.Error(err)
	}
	res.PrintModel(t.Log)
	if len(res.Classes) != 2 {
		t.Fatal("expected two classes", res.Classes)
	}
	stories := res.Classes[M.MakeStringId("stories")]
	if stories == nil {
		t.Fatal("expected stories", res.Classes)
	}
	if stories.Singular() != "story" {
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
	_, err := s.Compile(os.Stderr)
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
	res, err := s.Compile(os.Stderr)
	if err != nil {
		t.Fatal(err)
	}
	res.PrintModel(t.Log)
	objs := res.Classes[M.MakeStringId("kinds")]
	if objs == nil {
		t.Fatal("missing objs", res)
	}
	props := objs.Properties()
	if props == nil {
		t.Fatal("missing props", objs)
	}
	if len(props) != 2 {
		t.Fatal("unexpected props", props)
	}
	text := props[M.MakeStringId("text")]
	if text == nil {
		t.Fatal("missing text", props)
	}
	_, isText := text.(*M.TextProperty)
	if !isText {
		t.Errorf("unexpected property type %T", text)
	}
	num := props[M.MakeStringId("num")]
	if num == nil {
		t.Fatal("missing num", props)
	}
	_, isNum := num.(*M.NumProperty)
	if !isNum {
		t.Errorf("unexpected property type %T", num)
	}
}

//
// create an instance ( with no properties )
func TestInst(t *testing.T) {
	s := Script{}
	s.The("kind", Called("test"), Exists())
	res, err := s.Compile(os.Stderr)
	if err != nil {
		t.Error(err)
	}
	res.PrintModel(t.Log)
	if len(res.Instances) != 1 {
		t.Fatal("expected one instance", res.Instances)
	}
	test := res.Instances[M.MakeStringId("test")]
	if test == nil {
		t.Fatal("expected test instance", res.Instances)
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
	res, err := s.Compile(os.Stderr)
	if err != nil {
		t.Fatal(err)
	}
	res.PrintModel(t.Log)
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
	res, err := s.Compile(os.Stderr)
	if err != nil {
		t.Fatal(err)
	}
	res.PrintModel(t.Log)
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
	res, err := s.Compile(os.Stderr)
	if err != nil {
		t.Log(err)
	}
	//res.PrintModel(t.Log)
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

	res, err := s.Compile(os.Stderr)
	if err == nil {
		res.PrintModel(t.Log)
		t.Fatal("expected unscored story to report an issue")
	}
	t.Log("expected error:", err)
}
