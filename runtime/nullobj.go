package runtime

import (
	G "github.com/ionous/sashimi/game" // iobject
	"github.com/ionous/sashimi/util/ident"
	"log"
	"runtime" // go's runtime
)

// FIX? generate the implemenation via "go generate"
type NullObject struct {
	log *log.Logger
}

func (this *NullObject) println(args ...interface{}) {
	this.log.Println(append([]interface{}{"NullObject>"}, args...))
}

//
//
//
func (this *NullObject) String() string {
	return this.Name()
}

func (adapt *NullObject) Id() ident.Id {
	return ""
}

//
// FIX: maybe dont share the NullObject, and have unique names for them all
//
func (this *NullObject) Name() (ret string) {
	return "NullObject"
}

//
//
//
func (this *NullObject) Exists() bool {
	return false
}

//
//
//
func (this *NullObject) Class(cls string) bool {
	return false
}

//
//
//
func (this *NullObject) Is(c string) (ret bool) {
	pc, _, _, okay := runtime.Caller(0)
	if okay {
		name := runtime.FuncForPC(pc).Name()
		this.println(name, c)
	}
	return
}

//
//
//
func (this *NullObject) SetIs(c string) {
	pc, _, _, okay := runtime.Caller(0)
	if okay {
		name := runtime.FuncForPC(pc).Name()
		this.println(name, c)
	}
}

//
//
//
func (this *NullObject) Num(p string) (ret float32) {
	pc, _, _, okay := runtime.Caller(0)
	if okay {
		name := runtime.FuncForPC(pc).Name()
		this.println(name, p)
	}
	return
}

//
//
//
func (this *NullObject) SetNum(p string, v float32) {
	pc, _, _, okay := runtime.Caller(0)
	if okay {
		name := runtime.FuncForPC(pc).Name()
		this.println(name, p, v)
	}
}

//
//
//
func (this *NullObject) Object(c string) G.IObject {
	pc, _, _, okay := runtime.Caller(0)
	if okay {
		name := runtime.FuncForPC(pc).Name()
		this.println(name, c)
	}
	return this
}

//
//
//
func (this *NullObject) ObjectList(c string) (ret []G.IObject) {
	pc, _, _, okay := runtime.Caller(0)
	if okay {
		name := runtime.FuncForPC(pc).Name()
		this.println(name, c)
	}
	return
}

//
//
//
func (this *NullObject) SetObject(c string, _ G.IObject) {
	pc, _, _, okay := runtime.Caller(0)
	if okay {
		name := runtime.FuncForPC(pc).Name()
		this.println(name, c)
	}
}

//
//
//
func (this *NullObject) Text(p string) (ret string) {
	pc, _, _, okay := runtime.Caller(0)
	if okay {
		name := runtime.FuncForPC(pc).Name()
		this.println(name, p)
	}
	return
}

//
//
//
func (this *NullObject) SetText(p string, v string) {
	pc, _, _, okay := runtime.Caller(0)
	if okay {
		name := runtime.FuncForPC(pc).Name()
		this.println(name, p, v)
	}
}

//
//
//
func (this *NullObject) Go(s string, _ ...G.IObject) {
	pc, _, _, okay := runtime.Caller(0)
	if okay {
		name := runtime.FuncForPC(pc).Name()
		this.println(name, s)
	}
}

//
//
//
func (this *NullObject) Says(s string) {
	pc, _, _, okay := runtime.Caller(0)
	if okay {
		name := runtime.FuncForPC(pc).Name()
		this.println(name, s)
	}
}
