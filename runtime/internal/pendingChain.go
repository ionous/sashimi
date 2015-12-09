package internal

import G "github.com/ionous/sashimi/game"

// PendingChain implements runtime.IChain, it fills out the chained callback if the user links one in.
type PendingChain struct {
	node **ChainedCallback
	data *RuntimeAction // nouns, etc to assocuate with callbacks created via then.
}

func (c PendingChain) Then(cb G.Callback) {
	// check for empy pending chain structures
	if c.node != nil {
		*(c.node) = &ChainedCallback{c.data, cb}
	}
}
