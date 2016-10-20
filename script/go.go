package script

import G "github.com/ionous/sashimi/game"

// FIX: TESTING: add a function which returns a g.Go function? how to chain? etc.
func Go(phrase G.RuntimePhrase, phrases ...G.RuntimePhrase) G.OldCallback {
	return func(g G.Play) {
		g.Go(phrase, phrases...)
	}
}
