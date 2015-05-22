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

// type ClassDispatcher struct {
// 	bubbles, captures E.EventMap
// }
// //
// // Send the event to the class handleres for the passed cycle
// //
// func (this ClassDispatchers) DispatchClassEvent(evt E.IEvent, class *M.ClassInfo, captureCycle bool) (err error) {
// 	for cls := classs; err == nil && cls != nil; cls = cls.Parent() {
// 		if events, ok := dispatchers.getDispatcher(cls, captureCycle); ok {
// 			if e = events.HandleEvents(evt); e != nil {
// 				err = e
// 				break
// 			}
// 		}
// 	}
// 	return err
// }

// //
// // Retrieve the dispatcher for the passed key, only if it exists.
// //
// func (this ClassDispatchers) getDispatcher(cls *M.ClassInfo, captureCycle bool) (ret E.EventMap, okay bool) {
// 	if dispatcher, ok := this.all[key]; ok {
// 		if captureCycle {
// 			ret = dispatcher.captures
// 		} else {
// 			ret = dispatcher.bubbles
// 		}
// 		okay = len(ret) > 0
// 	}
// 	return ret, okay
// }
