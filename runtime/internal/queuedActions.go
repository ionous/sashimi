package internal

import (
	"fmt"
	E "github.com/ionous/sashimi/event"
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/util/ident"
)

// QueuedAction implements Future for named actions.
// ex. g.The("player").Go("jump")
type QueuedAction struct {
	data *RuntimeAction
	next *ChainedCallback
}

// QueuedAction implements Future for runtime phrases.
// ex. g.Go(Jump("player")
type QueuedPhrase struct {
	data *RuntimeAction
	run  G.RuntimePhrase
	next *ChainedCallback
}

func (c *QueuedPhrase) String() string {
	return fmt.Sprint("QueuedPhrase", c.run)
}

// QueuedPhrases implements Future for a set of phrases.
// One aspect not yet? modeled are long running actions.
// Currently, there would be no difference, therefore, between chaining on QueuedPhrases and chaining on the last element of QueuedPhrases
type QueuedPhrases struct {
	data *RuntimeAction
	run  []G.RuntimePhrase
	next *ChainedCallback
}

func (c *QueuedPhrases) String() string {
	return fmt.Sprint("QueuedPhrases", c.run)
}

// ChainedCallbacks implements Future for nested callbacks; created by PendingChain.
// ex. Then(func(g G.Play){ })
type ChainedCallback struct {
	data *RuntimeAction
	cb   G.Callback
}

func (c *ChainedCallback) String() string {
	return fmt.Sprint("ChainedCallback", c.cb)
}

// g.The("player").Go("hack", "the nice code").Then(trailing actions...)
func (a *QueuedAction) Run(g *Game) (err error) {
	// we've looped back now; end the event.
	act := a.data
	// start a new event frame:
	tgt := NewObjectTarget(g, act.GetTarget())
	path := E.NewPathTo(tgt)
	msg := &E.Message{Id: act.action.GetEvent().GetId(), Data: act}
	frame := g.Frame.BeginEvent(tgt, path, msg)
	// send the event, noting that new things may enter our queue.
	if runDefault, e := msg.Send(path); e != nil {
		err = e
	} else {
		// run default actions if requested, noting that new things may enter our queue.
		if !runDefault {
			frame.EndEvent()
		} else {
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
			}
			frame.EndEvent()

			// run "after" actions, which are queued dynamically ( though who knows why. )
			if after := a.data.after; len(after) > 0 {
				// fmt.Println(len(after), "after actions")
				play := g.newPlay(a.data, ident.Empty())
				for _, after := range after {
					after.call(play)
				}
			}
			// finally, run any trailing actions the caller may have specified.
			// this is done outside of the event frame, we will see these later...
			if a.next != nil {
				// fmt.Println("queuing then")
				a.next.Run(g)
			}
		}
	}
	return
}

func (c *QueuedPhrase) Run(g *Game) (err error) {
	if e := RunPhrases(g, c.data, c.run); e != nil {
		err = e
	} else if c.next != nil {
		g.Queue.QueueFuture(c.next)
	}
	return
}

//
func (c *QueuedPhrases) Run(g *Game) (err error) {
	if e := RunPhrases(g, c.data, c.run...); e != nil {
		err = e
	} else if c.next != nil {
		g.Queue.QueueFuture(c.next)
	}
	return
}

// during execute, whatever event we are in -- we want to stay in, until the end of execute.
// basically we want to subvert the queue -- so all queued futures come into us
func RunPhrases(g *Game, d *RuntimeAction, phrases ...G.RuntimePhrase) (err error) {
	// oldQueue := g.Queue
	// defer func() {
	// 	g.Queue = oldQueue
	// }()
	// myQueue := NewActionQueue()
	// // FIX: maybe we should get the queue from the event adapter...?
	// g.Queue = myQueue
	play := &GameEventAdapter{Game: g, data: d}
	for _, run := range phrases {
		run.Execute(play)
	}
	// return myQueue.ProcessActions(g)
	return
}

//
func (c *ChainedCallback) Run(g *Game) (err error) {
	play := &GameEventAdapter{Game: g, data: c.data}
	c.cb(play)
	return
}
