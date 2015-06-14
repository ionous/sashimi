package app

import (
	"github.com/ionous/sashimi/net/resource"
	"github.com/ionous/sashimi/net/session"
	"io"
)

const (
	frameKey = "frame"
)

//
// Transparently wraps a session resource, and all of its children.
//
type SessionResource struct {
	session  *CommandSession
	endpoint resource.IResource
}

//
// Helper to create a session and set its initial frame count.
//
func NewSessionResource(sessions *session.Sessions) resource.IResource {
	return &resource.Wrapper{
		Posts: func(_ io.Reader, doc resource.DocumentBuilder) (err error) {
			if _, sd, e := sessions.NewSession(); e != nil {
				err = e
			} else {
				session := sd.(*CommandSession)
				session.output.FlushDocument(doc)
				// now returning frames.
				doc.SetMeta(frameKey, session.frameCount)
			}
			return err
		}}
}

//
// Finds the sub-resource, and updates the internal endpoint
// Always returns "this".
//
func (this *SessionResource) Find(name string) (resource.IResource, bool) {
	child, okay := this.endpoint.Find(name)
	this.endpoint = child
	return this, okay
}

//
// Run a query on the endpoint, but do it inside a read lock
// And, add our turn-metadata to every document
//
func (this *SessionResource) Query() resource.Document {
	defer this.session.RUnlock()
	this.session.RLock()
	doc := this.endpoint.Query()
	resource.NewDocumentBuilder(&doc).SetMeta(frameKey, this.session.frameCount)
	return doc
}

//
// Run a post on the endpoint, but do it inside a write lock.
//
func (this *SessionResource) Post(reader io.Reader) (resource.Document, error) {
	// FIX? because the interface takes a reader,
	// we have to sit on the write lock for the duration of the read.
	defer this.session.Unlock()
	this.session.Lock()
	doc, err := this.endpoint.Post(reader)
	if err == nil {
		// by updating the frame first,
		// the response overrides any gets which may have happened just moments ago.
		this.session.frameCount++
		resource.NewDocumentBuilder(&doc).SetMeta(frameKey, this.session.frameCount)
	}
	return doc, err
}
