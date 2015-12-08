package native

import (
	G "github.com/ionous/sashimi/game"
)

// QuipHistory
// FIX: we just need fast sorting or queris; we've got recollection after all.
const QuipHistoryDepth = 3

type QuipHistory [QuipHistoryDepth]G.IValue

func (qh QuipHistory) Reset() {
	for _, h := range qh {
		h.SetObject(nil)
	}
}

func (qh QuipHistory) PushQuip(quip G.IObject) {
	for i := len(qh) - 1; i > 0; i = i - 1 {
		qh[i].SetObject(qh[i-1].Object())
	}
	qh[0].SetObject(quip)
}

func (qh QuipHistory) MostRecent() G.IObject {
	return qh[0].Object()
}

// Rank returns a number where larger is more recent, and 0 is not recent at all.
func (qh *QuipHistory) Rank(which G.IObject) (ret int) {
	for i, r := range qh {
		if r.Object() == which {
			ret = len(qh) - i
			break
		}
	}
	return ret
}
