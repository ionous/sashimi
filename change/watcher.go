package change

import (
	//"fmt"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
)

// FIX: add a generic shadow factory which creates instances, etc. from the api.
type ModelWatcher struct {
	PropertyChange
	meta.Model
}

func NewModelWatcher(ch PropertyChange, m meta.Model) meta.Model {
	return ModelWatcher{ch, m}
}

func (mw ModelWatcher) GetInstance(id ident.Id) (ret meta.Instance, okay bool) {
	//fmt.Println("get instance", id)
	if i, ok := mw.Model.GetInstance(id); ok {
		ret, okay = iwatch{mw, i}, true
	}
	return
}

func (mw ModelWatcher) InstanceNum(idx int) meta.Instance {
	//fmt.Println("instance num", idx)
	return iwatch{mw, mw.Model.InstanceNum(idx)}
}

type iwatch struct {
	mw ModelWatcher
	meta.Instance
}

func (iw iwatch) PropertyNum(idx int) meta.Property {
	//fmt.Println("property num", idx)
	p := iw.Instance.PropertyNum(idx)
	return pwatch{iw, p}
}

func (iw iwatch) FindProperty(name string) (ret meta.Property, okay bool) {
	//fmt.Println("find property", name)
	if p, ok := iw.Instance.FindProperty(name); ok {
		ret, okay = pwatch{iw, p}, true
	}
	return
}

func (iw iwatch) GetProperty(id ident.Id) (ret meta.Property, okay bool) {
	//fmt.Println("get property", id)
	if p, ok := iw.Instance.GetProperty(id); ok {
		ret, okay = pwatch{iw, p}, true
	}
	return
}
func (iw iwatch) GetPropertyByChoice(choice ident.Id) (ret meta.Property, okay bool) {
	//fmt.Println("get property by choice", choice)
	if p, ok := iw.Instance.GetPropertyByChoice(choice); ok {
		ret, okay = pwatch{iw, p}, true
	}
	return
}

type pwatch struct {
	iw iwatch
	meta.Property
}

func (pw pwatch) GetValue() meta.Value {
	//fmt.Println("get value", pw.GetId())
	return vwatch{pw, pw.Property.GetValue()}
}

// note: not currently watching arrays
func (pw pwatch) GetValues() meta.Values {
	//fmt.Println("get values", pw.GetId())
	return pw.Property.GetValues()
}

type vwatch struct {
	pw pwatch
	meta.Value
}

func (vw vwatch) SetNum(val float32) (err error) {
	//fmt.Println("set num", val)
	if old := vw.Value.GetNum(); old != val {
		if e := vw.Value.SetNum(val); e != nil {
			err = e
		} else {
			vw.pw.iw.mw.NumChange(
				vw.pw.iw.Instance,
				vw.pw.GetId(),
				old, val)
		}
	}
	return
}

func (vw vwatch) SetText(val string) (err error) {
	//fmt.Println("set text", val)
	if old := vw.Value.GetText(); old != val {
		if e := vw.Value.SetText(val); e != nil {
			err = e
		} else {
			vw.pw.iw.mw.TextChange(
				vw.pw.iw.Instance,
				vw.pw.GetId(),
				old, val)
		}
	}
	return
}
func (vw vwatch) SetState(val ident.Id) (err error) {
	//fmt.Println("set state", val)
	if old := vw.Value.GetState(); old != val {
		if e := vw.Value.SetState(val); e != nil {
			err = e
		} else {
			vw.pw.iw.mw.StateChange(
				vw.pw.iw.Instance,
				vw.pw.GetId(),
				old, val)
		}
	}
	return
}
func (vw vwatch) SetObject(val ident.Id) (err error) {
	//fmt.Println("set object", val)
	if old := vw.Value.GetObject(); old != val {
		if e := vw.Value.SetObject(val); e != nil {
			err = e
		} else {
			// notify
			var prev, next meta.Instance
			if !old.Empty() {
				prev, _ = vw.pw.iw.mw.GetInstance(old)
			}
			if !val.Empty() {
				next, _ = vw.pw.iw.mw.GetInstance(val)
			}
			// other is list of "contents", "inventory", etc.
			var other ident.Id
			if rel, ok := vw.pw.GetRelative(); ok {
				other = rel.From
			}
			//
			vw.pw.iw.mw.ReferenceChange(
				vw.pw.iw.Instance,
				vw.pw.GetId(),
				other,
				prev, next)
		}
	}
	return
}
