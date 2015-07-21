package extensions

import (
	"fmt"
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
	ret := quipScore{quip,
		quip.Text("reply") != "",
		quip.Is("repeatable"),
		false,
		100,
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
		if Debugging {
			fmt.Println(fmt.Sprintf(
				"%s replies:%t, repeats:%t, directly:%t, rank:%d",
				s.quip, s.replies, s.repeats, s.directly, s.rank))
		}
	}
	return ret
}

// QuipList returns the possible quips for the player to say.
func GetPlayerQuips(g G.Play) []G.IObject {
	qs := QuipSort{}
	con := g.Global("conversation").(*Conversation)
	if npc, ok := con.Interlocutor.Get(); ok {
		qh, qm := con.History, con.Memory
		latest := qh.MostRecent(g)
		isRestrictive := latest.Exists() && latest.Is("restrictive")

		// hrmm... this is very similar to "UpdateNextQuips"
		g.Visit("quips", func(quip G.IObject) bool {
			// Filter to quips which quip supply the interlocutor.
			if subject := quip.Object("subject"); subject == npc {
				// Filter to quips which have player comments.
				if quip.Text("comment") != "" {
					// Exclude one-time quips, checking the recollection table.
					if quip.Is("repeatable") || !qm.Recollects(quip) {
						// Check whether facts restrict this selection.
						if qm.IsQuipAllowed(g, quip) {
							// When following a restrictive quips, limit to those which directly follow.
							if isRestrictive && Quip(quip).Follows(latest).Directly(g) {
								qs.Add(quip)
							} else if !isRestrictive {
								// Select those which follow recent quips,
								// and those which do not have any follow constraints.
								if rank, direct := Quip(quip).Recently(qh).Follows(g); rank != 0 {
									score := qs.Add(quip)
									score.rank = rank
									score.directly = direct
								}
							}
						}
					}
				}
			}
			return false
		})
	}
	return qs.Sort()
}
