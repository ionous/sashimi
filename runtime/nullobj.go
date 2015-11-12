package runtime

import (
	"fmt"
	G "github.com/ionous/sashimi/game" // iobject
	"github.com/ionous/sashimi/util/ident"
	"log"
	"runtime"
)

// FIX? generate the implemenation via "go generate"
type _Null struct {
	name string
	pc   uintptr
}

func NullObject(name string) G.IObject {
	return NullObjectSource(name, 1)
}

func NullObjectSource(name string, skip int) G.IObject {
	pc := []uintptr{0}
	runtime.Callers(skip+1, pc) // 0 is callers itself, 1 is this code
	return _Null{name, pc[0]}
}

func (null _Null) println(args ...interface{}) {
	log.Println(append([]interface{}{null.String()}, args...))
}

//
func (null _Null) String() string {
	var str string
	f := runtime.FuncForPC(null.pc - 1)
	if f != nil {
		file, line := f.FileLine(null.pc - 1)
		str = fmt.Sprintf("(%s:%d)", file, line)
	}
	return "<NullObject " + null.name + str + ">"
}

//
func (null _Null) Id() ident.Id {
	return ""
}

//
func (null _Null) Exists() bool {
	return false
}

//
func (null _Null) FromClass(cls string) bool {
	return false
}

//
func (null _Null) Is(c string) (ret bool) {
	null.println("Is", c)
	return
}

//
func (null _Null) IsNow(c string) {
	null.println("IsNow", c)
}

//
func (null _Null) Get(p string) G.IValue {
	null.println("Get", p)
	return nullValue{fmt.Sprintf("%s.<%s>", null.String(), p)}
}

func (null _Null) GetList(p string) G.IList {
	null.println("GetList", p)
	return nullList{fmt.Sprintf("%s.<%s>", null.String(), p)}
}

//
func (null _Null) Num(p string) (ret float32) {
	null.println("Num", p)
	return
}

//
func (null _Null) SetNum(p string, v float32) {
	null.println("SetNum", p, v)
}

//
func (null _Null) Object(c string) G.IObject {
	null.println("Object", c)
	return null
}

//
func (null _Null) ObjectList(c string) (ret []G.IObject) {
	null.println("ObjectList", c)
	return
}

//
func (null _Null) Set(c string, _ G.IObject) {
	null.println("Set", c)
}

//
func (null _Null) Text(p string) (ret string) {
	null.println("Text", p)
	return
}

//
func (null _Null) SetText(p string, v string) {
	null.println("SetText", p, v)
}

//
func (null _Null) Go(s string, _ ...G.IObject) {
	null.println("Go", s)
}

//
func (null _Null) Says(s string) {
	null.println("Says", s)
}
