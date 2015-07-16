package extensions

import (
	G "github.com/ionous/sashimi/game"
	R "github.com/ionous/sashimi/runtime"
)

// FIX: we just need fast sorting.
// this is over-specification --
// we've got recollection after all.
type QuipHistory struct {
	r [QuipHistoryDepth]GosNilInterfacesAreAnnoying
}

const QuipHistoryDepth = 3

type GosNilInterfacesAreAnnoying struct {
	obj    G.IObject
	notnil bool
}

func (g GosNilInterfacesAreAnnoying) Get() (G.IObject, bool) {
	return g.obj, g.notnil && g.obj.Exists()
}
func (g *GosNilInterfacesAreAnnoying) Set(obj G.IObject) {
	g.obj = obj
	g.notnil = true
}
func (g *GosNilInterfacesAreAnnoying) Clear() {
	g.notnil = false
}

func (qh *QuipHistory) ClearQuips() {
	qh.r = QuipHistory{}.r
}
func (qh *QuipHistory) PushQuip(quip G.IObject) {
	qh.r[2], qh.r[1], qh.r[0] = qh.r[1], qh.r[0], GosNilInterfacesAreAnnoying{quip, true}
}
func (qh *QuipHistory) MostRecent(g G.Play) (andWhereAreTheTernaries G.IObject) {
	e := qh.r[0]
	if e.notnil {
		andWhereAreTheTernaries = e.obj
	} else {
		andWhereAreTheTernaries = R.NullObject("quip history")
	}
	return andWhereAreTheTernaries
}

// returns a rank where larger is more recent, and 0 is not recent at all.
func (qh *QuipHistory) Rank(which G.IObject) (ret int) {
	for i, r := range qh.r {
		if r.obj == which {
			ret = len(qh.r) - i
			break
		}
	}
	return ret
}
