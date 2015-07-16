package extensions

import (
	G "github.com/ionous/sashimi/game"
)

type quipMemories map[G.IObject]bool

var QuipMemory quipMemories = make(quipMemories)

// LearnQuip causes actors to recollect the passed quip.
// we can use this for facts for now too
// mostly the player will need this -- so just a table with precese is enough
// but it could also be actor, id
func (m quipMemories) Learn(quip G.IObject) {
	m[quip] = true
}

// RecollectsQuip determines if the passed quip has been spoken.
func (m quipMemories) Recollects(quip G.IObject) (recollects bool) {
	_, recollects = m[quip]
	return recollects
}
