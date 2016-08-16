package quips

import (
	"fmt"
	G "github.com/ionous/sashimi/game"
	"sort"
)

// QuipSort ranks quips by the importance of their comments.
type QuipSort struct {
	quips []quipScore
}

// record for tracking sorted scores
type quipScore struct {
	quip G.IObject
	rank int
}

func (s quipScore) String() string {
	return fmt.Sprintf("%s (%d)", s.quip, s.rank)
}

// Add the passed quip to the list of quips to sort
func (qs *QuipSort) Add(quip G.IObject) *quipScore {
	var rank int = 500
	if quip.Is("important") {
		rank = 1000
	} else if quip.Is("trivial") {
		rank = 200
	} else if quip.Is("departing") {
		rank = 100
	}
	if len(quip.Text("reply")) > 0 {
		rank += 50
	}
	ret := quipScore{quip, rank}
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

// sort.Interface; less is closer to the top
func (qs QuipSort) Less(i, j int) (moreImportant bool) {
	a, b := qs.quips[i], qs.quips[j]
	return a.rank > b.rank
}

func (qs QuipSort) Sort() []G.IObject {
	sort.Sort(qs)
	ret := make([]G.IObject, len(qs.quips))
	for i, s := range qs.quips {
		ret[i] = s.quip
	}
	return ret
}
