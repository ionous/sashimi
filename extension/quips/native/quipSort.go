package native

import (
	G "github.com/ionous/sashimi/game"
	"sort"
)

type QuipSort struct {
	quips []quipScore
}

// record for tracking sorted scores
type quipScore struct {
	quip                       G.IObject
	replies, repeats, directly bool
	rank                       int
}

// Add the passed quip to the list of quips to sort
func (qs *QuipSort) Add(quip G.IObject) *quipScore {
	var rank int = 250
	if quip.Is("important") {
		rank = 500
	} else if quip.Is("trivial") {
		rank = 100
	}
	ret := quipScore{quip,
		quip.Text("reply") != "",
		quip.Is("repeatable"),
		false,
		rank,
	}
	qs.quips = append(qs.quips, ret)
	return &qs.quips[len(qs.quips)-1]
}

// sort.Interface
func (qs QuipSort) Len() int {
	return len(qs.quips)
}

// sort.Interface
func (qs QuipSort) Swap(i, j int) {
	qs.quips[i], qs.quips[j] = qs.quips[j], qs.quips[i]
}

// sort.Interface; less is closer to the top of the lst
func (qs QuipSort) Less(i, j int) (less bool) {
	a, b := qs.quips[i], qs.quips[j]
	return (a.replies && !b.replies) ||
		(!a.repeats && b.repeats) ||
		(a.rank > b.rank) ||
		(a.directly && !b.directly)
}

func (qs QuipSort) Sort() []G.IObject {
	sort.Sort(qs)
	ret := make([]G.IObject, len(qs.quips))
	for i, s := range qs.quips {
		ret[i] = s.quip
		// if standard.Debugging {
		// 	g.Log(fmt.Sprintf(
		// 		"%s replies:%t, repeats:%t, directly:%t, rank:%d",
		// 		s.quip, s.replies, s.repeats, s.directly, s.rank))
		// }
	}
	return ret
}
