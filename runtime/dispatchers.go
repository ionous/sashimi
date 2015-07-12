package runtime

import (
	E "github.com/ionous/sashimi/event"
	"github.com/ionous/sashimi/util/ident"
	"log"
)

//
// Pool of all active class dispatchers
//
type Dispatchers struct {
	all map[ident.Id]Dispatcher
	log *log.Logger
}
type Dispatcher struct {
	E.Dispatcher
}

//
// Create a new dispatcher pool.
//
func NewDispatchers(log *log.Logger) Dispatchers {
	return Dispatchers{make(map[ident.Id]Dispatcher), log}
}

//
// Retrieve the dispatcher for the passed key, creating the dispatcher if it doesn't yet exist.
//
func (this Dispatchers) CreateDispatcher(id ident.Id) (ret Dispatcher) {
	if dispatcher, ok := this.all[id]; ok {
		ret = dispatcher
	} else {
		dispatcher.Dispatcher = E.NewDispatcher()
		this.all[id] = dispatcher
		ret = dispatcher
	}
	return ret
}

func (this Dispatchers) GetDispatcher(id ident.Id) (ret Dispatcher, okay bool) {
	ret, okay = this.all[id]
	return ret, okay
}
