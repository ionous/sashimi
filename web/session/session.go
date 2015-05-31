/*
 Give each new game its own set of data.
 Includes some handling for (unexpected) concurrent requests to the same sessions.
*/
package session

import (
	"fmt"
	"net/http"
)

//
// wrapper around ISession to handle concurrency.
//
type Session struct {
	id             string
	session        ISession
	inOut          chan CallResponse
	contentType    string
	valid, started bool
}

//
// helper to serialize concurrent requests.
//
type CallResponse struct {
	input  string
	output http.ResponseWriter
	done   Signal
}

//
// helper to block during call/response processing.
//
type Signal chan struct{}

//
// Send new input to the game, blocks until output has been written.
//
func (this *Session) Handle(in string, w http.ResponseWriter) {
	inOut, valid := this.inOut, this.valid
	if inOut != nil && valid {
		done := make(Signal)
		inOut <- CallResponse{in, w, done}
		<-done // block until the go routine has finished
	}
}

//
// start a background go routine to support potentially concurrent calls to Handle()
//
func (this *Session) Serve(contentType string) *Session {
	// FIX: some sort of timeout to evenually kill sessions?
	if this.inOut == nil {
		this.inOut, this.contentType, this.valid = make(chan CallResponse), contentType, true
		go func() {
			for {
				cr := <-this.inOut
				in, w, done := cr.input, cr.output, cr.done
				if e := this.handleInput(in, w); e != nil {
					if _, closed := e.(SessionClosed); closed {
						this.valid = false
						fmt.Println("!!! Closed")
						break
					} else {
						http.Error(w, e.Error(), http.StatusInternalServerError)
					}
				}
				done <- struct{}{}
			}
			fmt.Println("!!! Exiting")
		}()
	}
	return this
}

// the dance around "started" isnt very nice
// it has to do with the fact the initial post redirects to a get
// so we want data on that first get, but there's been no input.
func (this *Session) handleInput(in string, w http.ResponseWriter) (err error) {
	start := in == "start"
	if !this.started {
		if start {
			this.started = true
			fmt.Println("starting", in, "for", this.id)
			err = this.session.Write(w)
		}
	} else if start {
		fmt.Println("ignoring", in, "for", this.id)
	} else {
		fmt.Println("handling", in, "for", this.id)
		if this.contentType != "" {
			w.Header().Set("Content-Type", this.contentType)
		}
		err = this.session.Read(in).Write(w)
	}
	return err
}
