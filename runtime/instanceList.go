package runtime

import (
	"fmt"
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/util/ident"
)

type iList struct {
	game      *Game
	path      PropertyPath
	instances []api.Instance
}

func (l iList) AppendNum(float32)      {}
func (l iList) AppendText(string)      {}
func (l iList) AppendObject(G.IObject) {}
func (l iList) Reset()                 {}

func (l iList) Len() int {
	return len(l.instances)
}

func (l iList) Get(i int) (ret G.IValue) {
	if cnt := len(l.instances); i < 0 || i >= cnt {
		l.log(".Get(%d) out of range(%d).", i, cnt)
		ret = nullValue(l.path.InvalidIndex(i))
	} else {
		val := iasv{l.instances[i]}
		ret = gameValue{l.game, l.path.Index(i), api.ObjectProperty, val}
	}
	return
}

func (l iList) Contains(v interface{}) (yes bool) {
	// id := v.Object().Id()
	// for _, i := range l.instances {
	// 	if id == i.GetId() {
	// 		yes = true
	// 		break
	// 	}
	// }
	// return
	panic("not implemented")
}

func (l iList) log(format string, v ...interface{}) {
	suffix := fmt.Sprintf(format, v...)
	l.game.Println(l.path, suffix)
}

// ................

type iasv struct {
	api.Instance
}

func (i iasv) GetObject() ident.Id {
	return i.GetId()
}

func (iasv) GetNum() float32      { panic("not a number") }
func (iasv) SetNum(float32) error { panic("not a number") }

func (iasv) GetText() string      { panic("not a number") }
func (iasv) SetText(string) error { panic("not a number") }

func (iasv) GetState() ident.Id      { panic("not a number") }
func (iasv) SetState(ident.Id) error { panic("not a number") }

func (i iasv) SetObject(ident.Id) error {
	return fmt.Errorf("instance list cannot be changed")
}
