package ess

import (
	"github.com/ionous/sashimi/net/resource"
)

type ISessionFactory interface {
	NewSession(resource.DocumentBuilder) (ISession, error)
	GetSession(string) (ISession, bool)
}

type ISession interface {
	// RWLock the session to protect against a given client's browser reading/writing at the same time.
	RWLock
	// IResource starting point for all resources contained by the session.
	resource.IResource
	// Frame provides the current turn number so clients can reject gets returned from prior frames.
	// ( ex. a get followed by a post where the get returns slower than the post via a secondary browser connection. )
	Frame() int
}
