package web

import (
	"encoding/base64"
	S "github.com/ionous/sashimi/script"
	"github.com/ionous/sashimi/standard"
	"github.com/satori/go.uuid"
	//"os"
	"io/ioutil"
	"strings"
	"sync"
)

type Sessions struct {
	sessions map[string]*Session
	mutex    *sync.Mutex // for sessions
}

func NewSessions() Sessions {
	return Sessions{make(map[string]*Session), &sync.Mutex{}}
}

// return a new, game session session
func (this *Sessions) NewSession() (sess *Session, err error) {
	id := strings.TrimRight(base64.URLEncoding.EncodeToString(uuid.NewV4().Bytes()), "=")
	// FIX: it's very silly to have to init and compile each time.
	/**/
	if model, e := S.InitScripts().Compile(ioutil.Discard /*os.Stderr*/); e != nil {
		err = e
	} else {
		out := &SessionOutput{}
		if game, e := standard.NewStandardGame(model, out); e != nil {
			err = e
		} else if s, e := StartSession(id, game, out); e != nil {
			err = e
		} else {
			sess = s
			defer this.mutex.Unlock()
			this.mutex.Lock()
			this.sessions[id] = sess
		}
	}
	return sess, err
}

// return an existing session
func (this *Sessions) Session(id string) (*Session, bool) {
	defer this.mutex.Unlock()
	this.mutex.Lock()
	ret, ok := this.sessions[id]
	return ret, ok
}
