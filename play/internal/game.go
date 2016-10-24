package internal

import (
	E "github.com/ionous/sashimi/event"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
	"github.com/ionous/sashimi/util/sbuf"
)

type Game struct {
	PlayCore
	lastRandom int
}

func NewGame(core PlayCore) *Game {
	return &Game{core, -1}
}

func (g *Game) Log(args ...interface{}) {
	s := sbuf.New(args...).Line()
	g.Logger.Write([]byte(s))
}

func (g *Game) Random(exclusiveMax int) int {
	n := g.Rand.Intn(exclusiveMax)
	if n == g.lastRandom {
		n = (n + 1) % exclusiveMax
	}
	g.lastRandom = n
	return n
}

// g.The("player").Go("hack", "the nice code").Then(trailing actions...)
func (g *Game) RunAction(id ident.Id, params ...meta.Generic) (err error) {
	if act, e := g.Rtm.GetAction(id, params); e != nil {
		err = e
	} else {
		tgt, ctx := act.GetTarget(), act.GetContext()

		dispatch := &Dispatch{g, act.GetNouns(), params, nil}

		// begin and end event
		path := E.NewPathTo(ObjectTarget{dispatch, tgt})
		msg := &E.Message{Id: act.GetEvent().GetId(), Data: act}
		frame := g.Frame.BeginEvent(tgt, ctx, path, msg)
		frameEnded := false
		endFrame := func() {
			if !frameEnded {
				frame.EndEvent()
				frameEnded = true
			}
		}
		defer endFrame()

		// send the event
		if runDefault, e := msg.Send(path); e != nil {
			err = e
		} else if runDefault {
			// run the defaults if desired
			if e := act.RunDefault(); e != nil {
				err = e
			} else {
				endFrame()
				err = dispatch.RunAfterActions(act)
			}
		}
	}
	return
}
