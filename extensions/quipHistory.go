package extensions

import (
	"github.com/ionous/sashimi/util/ident"
)

const QuipHistoryDepth = 3

type QuipHistory struct{ r [QuipHistoryDepth]ident.Id }

func (qh *QuipHistory) Clear() {
	qh.r = QuipHistory{}.r
}
func (qh *QuipHistory) Push(id ident.Id) {
	qh.r[2], qh.r[1], qh.r[0] = qh.r[1], qh.r[0], id
}
func (qh *QuipHistory) MostRecent() (id ident.Id) {
	return qh.r[0]
}

// returns a rank where larger is more recent, and 0 is not recent at all.
func (qh *QuipHistory) Rank(id ident.Id) (ret int) {
	for i, r := range qh.r {
		if r == id {
			ret = len(qh.r) - i
			break
		}
	}
	return ret
}
