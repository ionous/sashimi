package internal

import (
	E "github.com/ionous/sashimi/event"
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/util/ident"
	"github.com/ionous/sashimi/util/sbuf"
	"log"
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
	return sbuf.New("QueuedPhrase:", c.run).String()
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
	return sbuf.New("QueuedPhrases:", c.run).String()
}

// ChainedCallbacks implements Future for nested callbacks; created by PendingChain.
// ex. Then(func(g G.Play){ })
type ChainedCallback struct {
	data *RuntimeAction
	cb   G.Callback
}

func (c *ChainedCallback) String() string {
	return sbuf.New("QueuedPhrases", c.cb).String()
}

// g.The("player").Go("hack", "the nice code").Then(trailing actions...)
func (a *QueuedAction) Run(g *Game) (err error) {
	log.Println("running queued action")
	// we've looped back now; end the event.
	act := a.data
	// start a new event frame:
	o1, o2 := act.GetTarget(), act.GetContext()
	path := E.NewPathTo(NewObjectTarget(g, o1))
	msg := &E.Message{Id: act.action.GetEvent().GetId(), Data: act}
	frame := g.Frame.BeginEvent(o1, o2, path, msg)
	// send the event, noting that new things may enter our queue.
	if runDefault, e := msg.Send(path); e != nil {
		err = e
	} else {
		// run default actions if requested, noting that new things may enter our queue.
		if !runDefault {
			frame.EndEvent()
		} else {
			if callbacks, ok := act.action.GetCallbacks(); ok {
				if cnt := callbacks.NumCallback(); cnt > 0 {
					rt := NewMars(g, g.newPlay(act, ident.Empty()))
					for i := 0; i < cnt; i++ {
						cb := callbacks.CallbackNum(i)
						if e := rt.Execute(cb); e != nil {
							err = e
							break
						}
					}
				}
			}
			frame.EndEvent()

			// run "after" actions, which are queued dynamically ( though who knows why. )
			if after := a.data.after; len(after) > 0 {
				// fmt.Println(len(after), "after actions")
				rt := NewMars(g, g.newPlay(a.data, ident.Empty()))
				for _, after := range after {
					if e := rt.Execute(after); e != nil {
						err = e
						break
					}
				}
			}
			// finally, run any trailing actions the caller may have specified.
			// this is done outside of the event frame, we will see these later...
			if a.next != nil {
				// fmt.Println("queuing then")
				err = a.next.Run(g)
			}
		}
	}
	return
}

func (c *QueuedPhrase) Run(g *Game) (err error) {
	log.Println("running queued phrase")
	if e := RunPhrases(g, c.data, c.run); e != nil {
		err = e
	} else if c.next != nil {
		g.Queue.QueueFuture(c.next)
	}
	return
}

//
func (c *QueuedPhrases) Run(g *Game) (err error) {
	log.Println("running queued phrases")
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
	play := &GameEventAdapter{Game: g, data: d}
	for _, run := range phrases {
		run.Execute(play)
	}
	return
}

//
func (c *ChainedCallback) Run(g *Game) (err error) {
	log.Println("running chained callback")
	rt := NewMars(g, &GameEventAdapter{Game: g, data: c.data})
	return rt.Execute(c.cb)
}
