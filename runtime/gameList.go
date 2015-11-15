package runtime

import (
	"fmt"
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/util/ident"
	"reflect"
)

var floatType = reflect.TypeOf(float64(0))

type gameList struct {
	game   *Game // context needed for wrapping instances
	path   PropertyPath
	ptype  api.PropertyType
	values api.Values
}

func (l gameList) Len() int {
	return l.values.NumValue()
}

func (l gameList) Get(i int) (ret G.IValue) {
	if cnt := l.values.NumValue(); i < 0 || i >= cnt {
		l.log("Get(%d) out of range(%d).", i, cnt)
		ret = nullValue{}
	} else {
		ret = gameValue{
			l.game,
			l.path.Index(i),
			l.ptype & ^api.ArrayProperty,
			l.values.ValueNum(i)}
	}
	return
}

func (l gameList) Contains(in interface{}) (yes bool) {
	switch l.ptype {
	default:
		panic("internal error, unhandled type")

	case api.NumProperty | api.ArrayProperty:
		if v, ok := in.(gameValue); ok && v.ptype == api.NumProperty {
			yes = containsNum(l.values, v.value.GetNum())
		} else {
			r := reflect.ValueOf(v)
			if r.Type().ConvertibleTo(floatType) {
				num := float32(r.Convert(floatType).Float())
				yes = containsNum(l.values, num)
			} else {
				l.log("Contains, testing unknown value %v", in)
			}
		}

	case api.TextProperty | api.ArrayProperty:
		if v, ok := in.(gameValue); ok && v.ptype == api.TextProperty {
			yes = containsText(l.values, v.value.GetText())
		} else if s, ok := in.(string); ok {
			yes = containsText(l.values, s)
		} else {
			l.log("Contains, testing unknown value %v", in)
		}

	case api.ObjectProperty | api.ArrayProperty:
		if in == nil {
			yes = containsObject(l.values, ident.Empty())
		} else if v, ok := in.(iasv); ok {
			yes = containsObject(l.values, v.GetId())
		} else if v, ok := in.(gameValue); ok && v.ptype == api.ObjectProperty {
			yes = containsObject(l.values, v.value.GetObject())
		} else if v, ok := in.(GameObject); ok {
			yes = containsObject(l.values, v.Id())
		} else if id, ok := in.(ident.Id); ok {
			yes = containsObject(l.values, id)
		} else {
			l.log("Contains, testing unknown value %v", in)
		}
	}
	return
}

func (l gameList) AppendNum(v float32) {
	if e := l.values.AppendNum(v); e != nil {
		l.log("AppendNum(): error appending list: %s.", e)
	}
}
func (l gameList) AppendText(v string) {
	if e := l.values.AppendText(v); e != nil {
		l.log("AppendText(): error appending list: %s.", e)
	}
}
func (l gameList) AppendObject(v G.IObject) {
	if e := l.values.AppendObject(v.Id()); e != nil {
		l.log("AppendObject(): error appending list: %s.", e)
	}
}
func (l gameList) Reset() {
	if e := l.values.ClearValues(); e != nil {
		l.log("Reset(): error reseting list: %s.", e)
	}
}

func (l gameList) log(format string, v ...interface{}) {
	suffix := fmt.Sprintf(format, v...)
	l.game.Println(l.path, suffix)
}

//.................

func containsNum(values api.Values, v float32) (yes bool) {
	for i := 0; i < values.NumValue(); i++ {
		if values.ValueNum(i).GetNum() == v {
			yes = true
			break
		}
	}
	return
}

func containsText(values api.Values, v string) (yes bool) {
	for i := 0; i < values.NumValue(); i++ {
		if values.ValueNum(i).GetText() == v {
			yes = true
			break
		}
	}
	return
}

func containsObject(values api.Values, v ident.Id) (yes bool) {
	for i := 0; i < values.NumValue(); i++ {
		if values.ValueNum(i).GetObject() == v {
			yes = true
			break
		}
	}
	return
}
