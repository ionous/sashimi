package runtime

import (
	E "github.com/ionous/sashimi/event"
	M "github.com/ionous/sashimi/model"
	"log"
)

//
// Pool of all active class dispatchers
//
type ClassDispatchers struct {
	all map[M.StringId]ClassDispatcher
	log *log.Logger
}
type ClassDispatcher struct {
	E.Dispatcher
}

//
// Create a new dispatcher pool.
//
func NewDispatchers(log *log.Logger) ClassDispatchers {
	return ClassDispatchers{make(map[M.StringId]ClassDispatcher), log}
}

//
// Retrieve the dispatcher for the passed key, creating the dispatcher if it doesn't yet exist.
//
func (this ClassDispatchers) CreateDispatcher(cls *M.ClassInfo) (ret ClassDispatcher) {
	if cls == nil {
		panic("nil passed to dispatcher creation")
	}
	id := cls.Id()
	if dispatcher, ok := this.all[id]; ok {
		ret = dispatcher
	} else {
		dispatcher.Dispatcher = E.NewDispatcher() //ClassDispatcher{make(E.EventMap), make(E.EventMap)}
		this.all[id] = dispatcher
		ret = dispatcher
	}
	return ret
}

func (this ClassDispatchers) GetDispatcher(cls *M.ClassInfo) (ret ClassDispatcher, okay bool) {
	if cls != nil {
		id := cls.Id()
		ret, okay = this.all[id]
	}
	return ret, okay
}
