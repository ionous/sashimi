package runtime

import (
	G "github.com/ionous/sashimi/game" // iobject
	"github.com/ionous/sashimi/util/ident"
	"log"
	"runtime" // go's runtime
)

// FIX? generate the implemenation via "go generate"
type _NullObject struct{}

func NullObject() G.IObject {
	return _NullObject{}
}

func (null _NullObject) println(args ...interface{}) {
	log.Println(append([]interface{}{"_NullObject>"}, args...))
}

//
func (null _NullObject) String() string {
	return "_NullObject"
}

//
func (null _NullObject) Id() ident.Id {
	return ""
}

//
func (null _NullObject) Exists() bool {
	return false
}

func (null _NullObject) Remove() {
}

//
func (null _NullObject) Class(cls string) bool {
	return false
}

//
func (null _NullObject) Is(c string) (ret bool) {
	pc, _, _, okay := runtime.Caller(0)
	if okay {
		name := runtime.FuncForPC(pc).Name()
		null.println(name, c)
	}
	return
}

//
func (null _NullObject) SetIs(c string) {
	pc, _, _, okay := runtime.Caller(0)
	if okay {
		name := runtime.FuncForPC(pc).Name()
		null.println(name, c)
	}
}

//
func (null _NullObject) Num(p string) (ret float32) {
	pc, _, _, okay := runtime.Caller(0)
	if okay {
		name := runtime.FuncForPC(pc).Name()
		null.println(name, p)
	}
	return
}

//
func (null _NullObject) SetNum(p string, v float32) {
	pc, _, _, okay := runtime.Caller(0)
	if okay {
		name := runtime.FuncForPC(pc).Name()
		null.println(name, p, v)
	}
}

//
func (null _NullObject) Object(c string) G.IObject {
	pc, _, _, okay := runtime.Caller(0)
	if okay {
		name := runtime.FuncForPC(pc).Name()
		null.println(name, c)
	}
	return null
}

//
func (null _NullObject) ObjectList(c string) (ret []G.IObject) {
	pc, _, _, okay := runtime.Caller(0)
	if okay {
		name := runtime.FuncForPC(pc).Name()
		null.println(name, c)
	}
	return
}

//
func (null _NullObject) Set(c string, _ G.IObject) {
	pc, _, _, okay := runtime.Caller(0)
	if okay {
		name := runtime.FuncForPC(pc).Name()
		null.println(name, c)
	}
}

//
func (null _NullObject) Text(p string) (ret string) {
	pc, _, _, okay := runtime.Caller(0)
	if okay {
		name := runtime.FuncForPC(pc).Name()
		null.println(name, p)
	}
	return
}

//
func (null _NullObject) SetText(p string, v string) {
	pc, _, _, okay := runtime.Caller(0)
	if okay {
		name := runtime.FuncForPC(pc).Name()
		null.println(name, p, v)
	}
}

//
func (null _NullObject) Go(s string, _ ...G.IObject) {
	pc, _, _, okay := runtime.Caller(0)
	if okay {
		name := runtime.FuncForPC(pc).Name()
		null.println(name, s)
	}
}

//
func (null _NullObject) Says(s string) {
	pc, _, _, okay := runtime.Caller(0)
	if okay {
		name := runtime.FuncForPC(pc).Name()
		null.println(name, s)
	}
}
