package live

// FIX: we have the concept floating in other fixes of "function" globals
// and that might be needed for this, where we really dont want *shared* globals
// you would want this tied to session, if at all possible.
//
const VersionString = "Sashimi Interactive Fiction Engine - 0.1"

var Debugging bool

func Debug() bool {
	return Debugging
}
