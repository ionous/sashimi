package internal

import (
	"github.com/ionous/mars/rt"
	"github.com/ionous/mars/rtm"
	E "github.com/ionous/sashimi/event"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/play/api"
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

func (g *Game) FindParent(obj rt.Object) (ret rt.Object, err error) {
	if p, e := g.Parents.LookupParent(obj.Instance); e != nil {
		err = e
	} else {
		ret = rt.Object{p}
	}
	return
}

//
func (g *Game) RunAction(id ident.Id, scp rt.Scope, args ...meta.Generic) (err error) {
	g.Log("act", id)
	if act, e := rtm.NewActionRuntime(g, id, scp, args); e != nil {
		err = e
	} else {
		// begin and end event
		tgt, ctx := act.GetTarget(), act.GetContext()
		// FIX? could this be an (object) stream? only if class was more uniform with instance
		// ( eg. a true prototype; also needed for default values )
		path := E.NewPathTo(ObjectTarget{act, tgt})
		msg := &E.Message{Id: act.GetRelatedEvent().GetId(), Data: act}
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
		runDefault, e := msg.Send(path)
		//
		if e != nil {
			err = e
		} else {
			if !runDefault {
				err = api.EventCancelled{}
			} else {
				if e := act.RunDefault(); e != nil {
					err = e
				} else {
					endFrame()
					err = act.RunAfterActions()
				}
			}
		}
	}
	return
}
