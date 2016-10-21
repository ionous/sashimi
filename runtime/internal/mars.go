package internal

import (
	"github.com/ionous/mars/rtm"
	"github.com/ionous/sashimi/meta"
)

type Mars struct {
	*rtm.Rtm
	ga *GameEventAdapter
}

func (rt *Mars) LookupParent(inst meta.Instance) (ret meta.Instance, rel meta.Property, okay bool) {
	return rt.ga.Game.LookupParent(inst)
}

func (rt *Mars) StopHere() {
	rt.ga.data.cancelled = true
}

func NewMars(g *Game, play *GameEventAdapter) *Mars {
	m := &Mars{rtm.NewRtm(g.Model), play}
	m.PushOutput(g.Writer)
	return m
}
