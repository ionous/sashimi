package session

import (
	"io"
)

// Generate a new session object for the passed id.
type SessionMaker func(id string) (ISession, error)

// Simple interface for running a game session.
// Input and output are separated into two phases to support game startup.
// Startup can generate output before the user submits their first command:
// for example, print banner text, describe rooms, etc.
type ISession interface {
	// Handle a single command input from the user.
	Read(string) ISession
	// Read the results of that command.
	Write(io.Writer) error
}
