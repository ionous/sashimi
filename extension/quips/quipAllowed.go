package quips

import (
	"github.com/ionous/sashimi/extension/facts"
	G "github.com/ionous/sashimi/game"
)

type QuipMemory struct {
	g G.Play
	facts.Memory
}

func PlayerMemory(g G.Play) QuipMemory {
	return QuipMemory{g, facts.PlayerMemory(g)}
}

func (qm QuipMemory) IsQuipAllowed(quip G.IObject) bool {
	return !qm.quipDisallowed(quip)
}

// IsQuipDisallowed evaluates the quip requirements and the known facts.
// A quip requirement can allow or disallow a given quip based on whether the player knows a specific fact.
func (qm QuipMemory) quipDisallowed(quip G.IObject) (disallowed bool) {
	// quip | fact | premits
	for reqs := qm.g.Query("quip requirements", true); reqs.HasNext(); {
		if req := reqs.Next(); quip.Equals(req.Get("quip").Object()) {
			fact, permits := req.Get("fact").Object(), req.Is("permitted")
			recollects := qm.Recollects(fact)
			// the opposite behavior of required is to exclude the use of the quip
			if permits != recollects {
				// the first disallowed quip/fact pairing stops the search because returning true stops.
				disallowed = true
				break
			}
		}
	}
	return disallowed
}
