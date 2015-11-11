package runtime

import (
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/util/ident"
)

// FIX: add some watcher tests?
// ( needs some sort of mock or model )
type ModelWatcher struct {
	api.Model
	PropertyChange
}

func NewModelWatcher(m api.Model, ch PropertyChange) api.Model {
	return ModelWatcher{m, ch}
}

func (mw ModelWatcher) GetInstance(id ident.Id) (ret api.Instance, okay bool) {
	if i, ok := mw.Model.GetInstance(id); ok {
		ret = iwatch{mw, i}
	}
	return
}

type iwatch struct {
	mw ModelWatcher
	api.Instance
}

func (iw iwatch) PropertyNum(i int) api.Property {
	p := iw.Instance.PropertyNum(i)
	return pwatch{iw, p}
}
func (iw iwatch) GetProperty(id ident.Id) (ret api.Property, okay bool) {
	if p, ok := iw.Instance.GetProperty(id); ok {
		ret = pwatch{iw, p}
		okay = true
	}
	return
}
func (iw iwatch) GetPropertyByChoice(choice ident.Id) (ret api.Property, okay bool) {
	if p, ok := iw.Instance.GetPropertyByChoice(choice); ok {
		ret = pwatch{iw, p}
		okay = true
	}
	return
}

type pwatch struct {
	iw iwatch
	api.Property
}

func (pw pwatch) GetValue() api.Value {
	return vwatch{pw, pw.Property.GetValue()}
}

func (pw pwatch) GetValues() api.Values {
	return zwatch{pw, pw.Property.GetValues()}
}

type vwatch struct {
	pw pwatch
	api.Value
}

func (vw vwatch) SetNum(val float32) (err error) {
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
	if old := vw.Value.GetObject(); old != val {
		if e := vw.Value.SetObject(val); e != nil {
			err = e
		} else {
			// notify
			var prev, next api.Instance
			if !old.Empty() {
				prev, _ = vw.pw.iw.mw.GetInstance(old)
			}
			if !val.Empty() {
				next, _ = vw.pw.iw.mw.GetInstance(val)
			}
			// other is looking for "contents", "inventory", etc.
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

type zwatch struct {
	pw pwatch
	api.Values
}

// func (zw zwatch) ClearValues() {
// }
// func (zw zwatch) AppendNum(float32) {
// }
// func (zw zwatch) AppendText(string) {
// }
// func (zw zwatch) AppendObject(ident.Id) error {
// }
