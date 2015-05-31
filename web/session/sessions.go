package session

import (
	"encoding/base64"
	"github.com/satori/go.uuid"
	"strings"
	"sync"
)

//
// wraps a factory to help handle concurrency.
//
type Sessions struct {
	factory     SessionMaker
	contentType string
	sessions    map[string]*Session
	mutex       *sync.Mutex // for sessions
}

//
// create a new session manager
//
func NewSessions(contentType string, factory SessionMaker) Sessions {
	return Sessions{factory, contentType, make(map[string]*Session), &sync.Mutex{}}
}

//
// Create a new game session, return its id.
//
func (this *Sessions) NewSession() (ret string, err error) {
	id := strings.TrimRight(base64.URLEncoding.EncodeToString(uuid.NewV4().Bytes()), "=")
	if sess, e := this.factory(); e != nil {
		err = e
	} else {
		s := (&Session{id: id, session: sess}).Serve(this.contentType)
		defer this.mutex.Unlock()
		this.mutex.Lock()
		this.sessions[id] = s
		ret = id
	}
	return ret, err
}

//
// return an existing session
//
func (this *Sessions) Session(id string) (*Session, bool) {
	defer this.mutex.Unlock()
	this.mutex.Lock()
	ret, ok := this.sessions[id]
	return ret, ok
}
