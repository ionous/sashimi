package runtime

import (
	G "github.com/ionous/sashimi/game" // iobject
	"github.com/ionous/sashimi/util/ident"
	"log"
	"runtime" // go's runtime
)

// FIX? generate the implemenation via "go generate"
type NullObject struct{}

func (null NullObject) println(args ...interface{}) {
	log.Println(append([]interface{}{"NullObject>"}, args...))
}

//
func (null NullObject) String() string {
	return "NullObject"
}

//
func (null NullObject) Id() ident.Id {
	return ""
}

//
func (null NullObject) Exists() bool {
	return false
}

//
func (null NullObject) Class(cls string) bool {
	return false
}

//
func (null NullObject) Is(c string) (ret bool) {
	pc, _, _, okay := runtime.Caller(0)
	if okay {
		name := runtime.FuncForPC(pc).Name()
		null.println(name, c)
	}
	return
}

//
func (null NullObject) SetIs(c string) {
	pc, _, _, okay := runtime.Caller(0)
	if okay {
		name := runtime.FuncForPC(pc).Name()
		null.println(name, c)
	}
}

//
func (null NullObject) Num(p string) (ret float32) {
	pc, _, _, okay := runtime.Caller(0)
	if okay {
		name := runtime.FuncForPC(pc).Name()
		null.println(name, p)
	}
	return
}

//
func (null NullObject) SetNum(p string, v float32) {
	pc, _, _, okay := runtime.Caller(0)
	if okay {
		name := runtime.FuncForPC(pc).Name()
		null.println(name, p, v)
	}
}

//
func (null NullObject) Object(c string) G.IObject {
	pc, _, _, okay := runtime.Caller(0)
	if okay {
		name := runtime.FuncForPC(pc).Name()
		null.println(name, c)
	}
	return null
}

//
func (null NullObject) ObjectList(c string) (ret []G.IObject) {
	pc, _, _, okay := runtime.Caller(0)
	if okay {
		name := runtime.FuncForPC(pc).Name()
		null.println(name, c)
	}
	return
}

//
func (null NullObject) Set(c string, _ G.IObject) {
	pc, _, _, okay := runtime.Caller(0)
	if okay {
		name := runtime.FuncForPC(pc).Name()
		null.println(name, c)
	}
}

//
func (null NullObject) Text(p string) (ret string) {
	pc, _, _, okay := runtime.Caller(0)
	if okay {
		name := runtime.FuncForPC(pc).Name()
		null.println(name, p)
	}
	return
}

//
func (null NullObject) SetText(p string, v string) {
	pc, _, _, okay := runtime.Caller(0)
	if okay {
		name := runtime.FuncForPC(pc).Name()
		null.println(name, p, v)
	}
}

//
func (null NullObject) Go(s string, _ ...G.IObject) {
	pc, _, _, okay := runtime.Caller(0)
	if okay {
		name := runtime.FuncForPC(pc).Name()
		null.println(name, s)
	}
}

//
func (null NullObject) Says(s string) {
	pc, _, _, okay := runtime.Caller(0)
	if okay {
		name := runtime.FuncForPC(pc).Name()
		null.println(name, s)
	}
}
