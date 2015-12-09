package internal

import (
	"container/list"
	E "github.com/ionous/sashimi/event"
	G "github.com/ionous/sashimi/game"
)

type ActionQueue struct {
	list *list.List
}

// FUTURE: do we need to pass Game here? could we have a future queue ( and promises ) as a completely separate system?
type Future interface {
	Run(*Game) error
}

//
func NewActionQueue() ActionQueue {
	return ActionQueue{list.New()}
}

//phrase G.RuntimePhrase
func (q ActionQueue) QueueFuture(f Future) {
	q.list.PushBack(f)
}

//
func (q ActionQueue) Empty() bool {
	return q.list.Front() == nil
}

//
func (q ActionQueue) PopFuture() (ret Future) {
	if el := q.list.Front(); el != nil {
		ret = q.list.Remove(el).(Future)
	}
	return ret
}

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
	// FIX: this should be the same as "defaultAction(s)"
	run  G.RuntimePhrase
	next *ChainedCallback
}

// QueuedPhrases implements Future for a set of phrases.
// One aspect not yet? modeled are long running actions.
// Currently, there would be no difference, therefore, between chaining on QueuedPhrases and chaining on the last element of QueuedPhrases
type QueuedPhrases struct {
	data *RuntimeAction
	run  []G.RuntimePhrase
	next *ChainedCallback
}

// ChainedCallbacks implements Future for nested callbacks; created by PendingChain.
// ex. Then(func(g G.Play){ })
type ChainedCallback struct {
	data *RuntimeAction
	cb   G.Callback
}

//
func (c *QueuedAction) Run(g *Game) (err error) {
	act := c.data
	tgt := NewObjectTarget(g, act.GetTarget())
	path := E.NewPathTo(tgt)
	msg := &E.Message{Id: act.action.GetEvent().GetId(), Data: act}
	frame := g.Frame.BeginEvent(tgt, path, msg)
	if runDefault, e := msg.Send(path); e != nil {
		err = e
	} else {
		if runDefault {
			err = act.runDefaultActions()
		}
	}
	frame.EndEvent()
	if err == nil && c.next != nil {
		g.Queue.QueueFuture(c.next)
	}
	return
}

//
func (c *QueuedPhrase) Run(g *Game) (err error) {
	play := &GameEventAdapter{Game: g, data: c.data}
	c.run.Execute(play)
	if c.next != nil {
		g.Queue.QueueFuture(c.next)
	}
	return
}

//
func (c *QueuedPhrases) Run(g *Game) (err error) {
	play := &GameEventAdapter{Game: g, data: c.data}
	for _, run := range c.run {
		run.Execute(play)
	}
	if c.next != nil {
		g.Queue.QueueFuture(c.next)
	}
	return
}

//
func (c *ChainedCallback) Run(g *Game) (err error) {
	play := &GameEventAdapter{Game: g, data: c.data}
	c.cb(play)
	return
}
