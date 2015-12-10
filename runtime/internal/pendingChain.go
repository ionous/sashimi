package internal

import G "github.com/ionous/sashimi/game"

// PendingChain implements runtime.IChain, it fills out the chained callback if the user links one in.
type PendingChain struct {
	//node **ChainedCallback
	src *GameEventAdapter
}

func NewPendingChain(src *GameEventAdapter, _ Future) PendingChain {
	//return PendingChain{src, future.GetChain()}
	return PendingChain{src}
}

func (c PendingChain) Then(cb G.Callback) {
	// watch for null objects
	if c.src != nil {
		// if c.node != nil {
		// 	*(c.node) = &ChainedCallback{c.data, cb}
		// }
		chain := ChainedCallback{c.src.data, cb}
		chain.Run(c.src.Game)
	}
}
