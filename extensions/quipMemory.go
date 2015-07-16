package extensions

import (
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
)

type QuipMemory struct {
	mem map[G.IObject]bool
}

// LearnQuip causes actors to recollect the passed quip.
// we can use this for facts for now too
// mostly the player will need this -- so just a table with precese is enough
// but it could also be actor, id
func (m *QuipMemory) Learn(quip G.IObject) {
	if m.mem == nil {
		m.mem = make(map[G.IObject]bool)
	}
	m.mem[quip] = true
}

// RecollectsQuip determines if the passed quip has been spoken.
func (m QuipMemory) Recollects(quip G.IObject) (recollects bool) {
	if m.mem != nil {
		_, recollects = m.mem[quip]
	}
	return recollects
}

func (m QuipMemory) IsQuipAllowed(g G.Play, quip G.IObject) bool {
	return !m.IsQuipDisallowed(g, quip)
}

// IsQuipDisallowed evaluates the quip requirements and the known facts.
// A quip requirement can allow or disallow a given quip based on whether the player knows a specific fact.
func (m QuipMemory) IsQuipDisallowed(g G.Play, quip G.IObject) bool {
	disallowed := g.Visit("quip requirements",
		func(req G.IObject) (disallowed bool) {
			if req.Object("quip") == quip {
				fact, permits := req.Object("fact"), req.Is("permitted")
				recollects := m.Recollects(fact)
				// the opposite behavior of required is to exclude the use of the quip
				if permits != recollects {
					disallowed = true
				}
			}
			// the first disallowed quip/fact pairing stops the search because returning true stops.
			return disallowed
		})
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
