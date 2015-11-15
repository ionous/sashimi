package extensions

import (
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
)

type QuipMemory struct {
	G.IList
}

func PlayerMemory(g G.Play) QuipMemory {
	return QuipMemory{g.The("player").List("recollections")}
}

func (mem QuipMemory) TriesToLearn(quip G.IObject) (newlyLearned bool) {
	if recollects := mem.Recollects(quip); !recollects {
		mem.Learns(quip)
		newlyLearned = true
	}
	return newlyLearned
}

// LearnQuip causes actors to recollect the passed quip.
// we can use this for facts for now too
// mostly the player will need this -- so just a table with precese is enough
// but it could also be actor, id
func (mem QuipMemory) Learns(quip G.IObject) {
	mem.AppendObject(quip)
}

// RecollectsQuip determines if the passed quip has been spoken.
func (mem QuipMemory) Recollects(quip G.IObject) bool {
	return mem.Contains(quip)
}

func (mem QuipMemory) IsQuipAllowed(g G.Play, quip G.IObject) bool {
	return !mem.IsQuipDisallowed(g, quip)
}

// IsQuipDisallowed evaluates the quip requirements and the known facts.
// A quip requirement can allow or disallow a given quip based on whether the player knows a specific fact.
func (mem QuipMemory) IsQuipDisallowed(g G.Play, quip G.IObject) (disallowed bool) {
	for i, reqs := 0, g.List("quip requirements"); i < reqs.Len(); i++ {
		if req := reqs.Get(i).Object(); req.Get("quip").Object() == quip {
			fact, permits := req.Get("fact").Object(), req.Is("permitted")
			recollects := mem.Recollects(fact)
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

// The("quip requirements", Table("fact", "permitted-property", "quip")...)
func init() {
	AddScript(func(s *Script) {
		s.The("kinds", Called("facts"),
			// FIX: interestingly, kinds should have names
			// having the same property as a parent class probably shouldnt be an error
			Have("summary", "text"))

		// FIX: should be "data"
		s.The("kinds", Called("quip requirements"),
			Have("fact", "fact"),
			AreEither("permitted").Or("prohibited"),
			Have("quip", "quip"))
	})
}
