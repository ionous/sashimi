package ess

// RWLock because go's built in "sync" doesnt expose interfaces.
type RWLock interface {
	RLock()
	RUnlock()
	Lock()
	Unlock()
}
