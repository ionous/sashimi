package extensions

import (
	G "github.com/ionous/sashimi/game"
)

const QuipHistoryDepth = 3

type GosNilInterfacesAreAnnoying struct {
	obj    G.IObject
	notnil bool
}

type QuipHistory struct {
	r [QuipHistoryDepth]GosNilInterfacesAreAnnoying
}

func (qh *QuipHistory) ClearQuips() {
	qh.r = QuipHistory{}.r
}
func (qh *QuipHistory) PushQuip(next G.IObject) {
	qh.r[2], qh.r[1], qh.r[0] = qh.r[1], qh.r[0], GosNilInterfacesAreAnnoying{next, true}
}
func (qh *QuipHistory) MostRecent(g G.Play) (andWhereAreTheTernaries G.IObject) {
	e := qh.r[0]
	if e.notnil {
		andWhereAreTheTernaries = e.obj
	} else {
		andWhereAreTheTernaries = g.The("")
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
