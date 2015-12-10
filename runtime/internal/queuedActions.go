package internal

import (
	"fmt"
	G "github.com/ionous/sashimi/game"
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

// BUG: we may be getting two sets of events?
// when we run the defaults, they may bqueue
// they will ahve the act of the parent
func (c *QueuedAction) Run(g *Game) (err error) {
	panic("not implemented")
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
