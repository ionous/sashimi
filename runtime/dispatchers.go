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
// CreateDispatcher, or retrieve one if it already exists
func (d Dispatchers) CreateDispatcher(id ident.Id) (ret Dispatcher) {
	if dispatcher, ok := d.all[id]; ok {
		ret = dispatcher
	} else {
		dispatcher.Dispatcher = E.NewDispatcher()
		d.all[id] = dispatcher
		ret = dispatcher
	}
	return ret
}

func (d Dispatchers) GetDispatcher(id ident.Id) (ret Dispatcher, okay bool) {
	ret, okay = d.all[id]
	return ret, okay
}
