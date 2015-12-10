package internal

import (
	"container/list"
	"fmt"
	E "github.com/ionous/sashimi/event"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/util/ident"
)

type ActionQueue struct {
	list *list.List
}

// FUTURE: do we need to pass Game here? could we have a future queue ( and promises ) as a completely separate system?
type Future interface {
	Run(g *Game) error
}

//
func NewActionQueue() ActionQueue {
	return ActionQueue{list.New()}
}

//
func (q ActionQueue) Empty() bool {
	return q.list.Front() == nil
}

//
func (q ActionQueue) QueueFuture(f Future) {
	q.list.PushBack(f)
	/**/ fmt.Println(fmt.Sprintf("queuing %T future, %d", f, q.list.Len()))
}

//
func (q ActionQueue) Pop() (ret Future) {
	if el := q.list.Front(); el != nil {
		ret = q.list.Remove(el).(Future)
	}
	/**/ fmt.Println(fmt.Sprintf("popped %T future, %d", ret, q.list.Len()))
	return ret
}

// this craziness exists to help unwind the very deep callstacks the events create
func (q *ActionQueue) ProcessActions(g *Game) (err error) {
	type PendingAction struct {
		a *QueuedAction
		api.IEndEvent
		list *list.List
	}
	pending := []PendingAction{}
	//
	for {
		if q.Empty() {
			if p := len(pending) - 1; p < 0 {
				/**/ fmt.Println("done")
				break // done, done, done.
			} else {
				end := pending[p]
				pending = pending[0:p]
				//
				/**/ fmt.Println("ending", end.a.data.action.GetId())
				end.EndEvent()
				q.list = end.list
				// run "after" actions, which are queued dynamically ( though who knows why. )
				if after := end.a.data.after; len(after) > 0 {
					/**/ fmt.Println(len(after), "after actions")
					play := g.newPlay(end.a.data, ident.Empty())
					for _, after := range after {
						after.call(play)
					}
				}
				// finally, run any trailing actions the caller may have specified.
				// this is done outside of the event frame, we will see these later...
				if end.a.next != nil {
					/**/ fmt.Println("queuing then")
					q.QueueFuture(end.a.next)
				}
			}
		} else {
			// normal action?
			r := q.Pop()
			if a, ok := r.(*QueuedAction); !ok {
				/**/ fmt.Println("running some other action", r)
				if e := r.Run(g); e != nil {
					err = e
					break
				}
			} else {
				// handling an action of the form:
				// g.The("player").Go("hack", "the nice code").Then(trailing actions...)
				// we've looped back now; end the event.
				act := a.data
				// start a new event frame:
				tgt := NewObjectTarget(g, act.GetTarget())
				path := E.NewPathTo(tgt)
				msg := &E.Message{Id: act.action.GetEvent().GetId(), Data: act}
				frame := g.Frame.BeginEvent(tgt, path, msg)
				// record this in our pseudo-stack
				pending = append(pending, PendingAction{a, frame, q.list})
				q.list = list.New()
				//
				// send the event, noting that new things may enter our queue.
				runDefault, e := msg.Send(path)

				if e != nil {
					err = e
					break
				} else {
					// run default actions if requested, noting that new things may enter our queue.
					if runDefault {
						play := g.newPlay(act, ident.Empty())
						if callbacks, ok := act.action.GetCallbacks(); ok {
							for i := 0; i < callbacks.NumCallback(); i++ {
								cb := callbacks.CallbackNum(i)
								if found, ok := g.LookupCallback(cb); ok {
									found(play)
								} else {
									err = fmt.Errorf("internal error, couldnt find callback %s", cb)
									break
								}
							}
							if err != nil {
								break
							}
						}
					}
				}
			}
		}
	}
	return err
}
