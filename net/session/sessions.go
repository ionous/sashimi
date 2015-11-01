package session

import (
	"encoding/base64"
	"github.com/satori/go.uuid"
	"strings"
	"sync"
)

//
// User session data.
//
type SessionData interface{}

//
// Generate a new session object for the passed id.
//
type SessionMaker func(id string) (SessionData, error)

//
// Create a new session manager.
//
func NewSessions(factory SessionMaker) *Sessions {
	return &Sessions{factory, make(map[string]SessionData), &sync.Mutex{}}
}

//
// Session data manager.
//
type Sessions struct {
	factory  SessionMaker
	sessions map[string]SessionData
	mutex    *sync.Mutex // for sessions
}

func (s *Sessions) NumSessions() int {
	defer s.mutex.Unlock()
	s.mutex.Lock()
	return len(s.sessions)
}

//
// Create a new game session, return its id.
//
func (s *Sessions) NewSession() (newId string, newData SessionData, err error) {
	id := strings.TrimRight(base64.URLEncoding.EncodeToString(uuid.NewV4().Bytes()), "=")
	if sessionData, e := s.factory(id); e != nil {
		err = e
	} else {
		//s := //(&Session{id: id, session: sess}).Serve()
		defer s.mutex.Unlock()
		s.mutex.Lock()
		s.sessions[id] = sessionData
		newId, newData = id, sessionData
	}
	return newId, newData, err
}

//
// Return an existing session.
//
func (s *Sessions) Session(id string) (ret SessionData, okay bool) {
	if id != "" {
		defer s.mutex.Unlock()
		s.mutex.Lock()
		ret, okay = s.sessions[id]
	}
	return ret, okay
}
