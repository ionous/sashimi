package runtime

import (
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/util/ident"
)

type propValue struct {
	oa    ObjectAdapter
	prop  api.Property
	value api.Value
}

func (n propValue) Num() (ret float32) {
	if t := n.prop.GetType(); t != api.NumProperty {
		n.oa.log("%s.Num(%s): property type(%d) is not a number.", n.oa, n.prop, t)
	} else {
		ret = n.value.GetNum()
	}
	return
}

func (n propValue) SetNum(value float32) {
	if t := n.prop.GetType(); t != api.NumProperty {
		n.oa.log("%s.SetNum(%s): property type(%d) is not a number.", n.oa, n.prop, t)
	} else if e := n.value.SetNum(value); e != nil {
		n.oa.log("%s.SetNum(%s): error setting value: %s.", n.oa, n.prop, e)
	}
}

func (n propValue) Text() (ret string) {
	if t := n.prop.GetType(); t != api.TextProperty {
		n.oa.log("%s.Text(%s): property type(%d) is not text.", n.oa, n.prop, t)
	} else {
		ret = n.value.GetText()
	}
	return
}

func (n propValue) SetText(text string) {
	if t := n.prop.GetType(); t != api.TextProperty {
		n.oa.log("%s.SetText(%s): property type(%d) is not text.", n.oa, n.prop, t)
	} else if e := n.value.SetText(text); e != nil {
		n.oa.log("%s.SetText(%s): error setting value: %s.", n.oa, n.prop, e)
	}
}

// TBD: should these be logged? its sure nice to have be able to test objects generically for properties
func (n propValue) Object() G.IObject {
	var res ident.Id
	if t := n.prop.GetType(); t != api.ObjectProperty {
		n.oa.log("%s.Object(%s): property type(%d) is not an object.", n.oa, n.prop, t)
	} else {
		res = n.prop.GetValue().GetObject()
	}
	return NewObjectAdapterFromId(n.oa.game, res)
}

func (n propValue) SetObject(object G.IObject) {
	switch t := n.prop.GetType(); t {
	default:
		n.oa.log("%s.Set(%s): property type(%d) is not an object.", n.oa, n.prop, t)

	case api.ObjectProperty:
		var id ident.Id
		if other, ok := object.(ObjectAdapter); ok {
			id = other.gobj.GetId()
		}
		if e := n.prop.GetValue().SetObject(id); e != nil {
			n.oa.log("%s.Set(%s): error setting value: %s.", n.oa, n.prop, e)
		}

	case api.ObjectProperty | api.ArrayProperty:
		values := n.prop.GetValues()
		if other, ok := object.(ObjectAdapter); !ok {
			if e := values.ClearValues(); e != nil {
				n.oa.log("%s.Set(%s): error clearing value: %s.", n.oa, n.prop, e)
			}
		} else {
			if e := values.AppendObject(other.gobj.GetId()); e != nil {
				n.oa.log("%s.Set(%s): error appending value: %s.", n.oa, n.prop, e)
			}
		}
	}

}
