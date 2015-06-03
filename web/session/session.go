package session

import (
	"fmt"
	"log"
)

//
// wrapper around ISession to handle concurrency.
//
type Session struct {
	id             string
	session        ISession
	inOut          chan CallResponse
	valid, started bool
}

// FIXIXIX remove
func (this *Session) Session() ISession {
	return this.session
}

//
// helper to serialize concurrent requests.
//
type CallResponse struct {
	input  interface{}
	output chan Output
}

//
//
//
type Response func(data interface{}) error

//
//
//
type Error func(error) bool

type Output struct {
	data interface{}
	err  error
}

//
// Send new input to the game, blocks until output has been written.
// The error handlers is called when a session returns an error.
// If the error handler returns false, then the background routine will exit.
//
func (this *Session) WriteRead(in interface{}) (ret interface{}, err error) {
	if inOut := this.inOut; inOut == nil {
		err = fmt.Errorf("session not started")
	} else if valid := this.valid; !valid {
		err = fmt.Errorf("session closed")
	} else {
		// FIX? if Write() returned the interface to Read() then session could transparently implement ISession
		done := make(chan Output)
		inOut <- CallResponse{in, done}
		output := <-done // block until the go routine has finished
		ret, err = output.data, output.err
	}
	return
}

//
// Start a background go routine to support potentially concurrent calls to Handle().
func (this *Session) Serve() *Session {
	// FIX: some sort of timeout to evenually kill sessions?
	if this.inOut == nil {
		this.inOut, this.valid = make(chan CallResponse), true
		go func() {
			for this.valid {
				cr := <-this.inOut
				data, err := this.handleInput(cr.input)
				if _, closed := err.(SessionClosed); closed {
					this.valid = false
				}
				cr.output <- Output{data, err}
			}
		}()
	}
	return this
}

// the dance around "started" isnt very nice
// it has to do with the fact the initial post redirects to a get
// so we want data on that first get, but there's been no input.
func (this *Session) handleInput(in interface{}) (ret interface{}, err error) {
	start := in == "start"
	if !this.started {
		if start {
			this.started = true
			log.Println("starting", this.id)
			ret, err = this.session.Read()
		}
	} else if start {
		log.Println("ignoring", in, "for", this.id)
	} else {
		log.Println("handling", in, "for", this.id)
		ret, err = this.session.Write(in).Read()
	}
	return ret, err
}
