package session

import (
	"encoding/base64"
	"github.com/satori/go.uuid"
	"strings"
	"sync"
)

//
// Session manager.
//
type Sessions struct {
	factory  SessionMaker
	sessions map[string]*Session
	mutex    *sync.Mutex // for sessions
}

//
// Create a new session manager.
//
func NewSessions(factory SessionMaker) Sessions {
	return Sessions{factory, make(map[string]*Session), &sync.Mutex{}}
}

//
// Create a new game session, return its id.
//
func (this *Sessions) NewSession() (newId string, err error) {
	id := strings.TrimRight(base64.URLEncoding.EncodeToString(uuid.NewV4().Bytes()), "=")
	if sess, e := this.factory(id); e != nil {
		err = e
	} else {
		s := (&Session{id: id, session: sess}).Serve()
		defer this.mutex.Unlock()
		this.mutex.Lock()
		this.sessions[id] = s
		newId = id
	}
	return newId, err
}

//
// Return an existing session.
//
func (this *Sessions) Session(id string) (ret *Session, okay bool) {
	if id != "" {
		defer this.mutex.Unlock()
		this.mutex.Lock()
		ret, okay = this.sessions[id]
	}
	return ret, okay
}
