package session

import (
	"io"
)

type SessionMaker func() (ISession, error)

type ISession interface {
	Read(string) ISession
	Write(io.Writer) error
}
