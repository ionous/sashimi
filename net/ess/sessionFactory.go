package ess

import (
	"github.com/ionous/sashimi/net/resource"
)

type SessionFactory interface {
	NewSession(resource.DocumentBuilder) (Session, error)
	GetSession(string) (Session, bool)
}
