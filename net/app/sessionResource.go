package app

import (
	"github.com/ionous/sashimi/net/resource"
	"github.com/ionous/sashimi/net/session"
	"io"
)

const (
	frameKey = "frame"
)

// SessionResource wraps a session as a resource, along with its children.
type SessionResource struct {
	session  *CommandSession
	endpoint resource.IResource
}

func NewSessionResource(sessions *session.Sessions) resource.IResource {
	return &resource.Wrapper{
		Posts: func(_ io.Reader, doc resource.DocumentBuilder) (err error) {
			if _, sd, e := sessions.NewSession(); e != nil {
				err = e
			} else {
				session := sd.(*CommandSession)
				session.out.FlushDocument(doc)
				// now returning frames.
				doc.SetMeta(frameKey, session.FrameCount())
			}
			return err
		}}
}

// Find the sub-resource, and updates the internal endpoint
// Always returns "res".
func (res *SessionResource) Find(name string) (resource.IResource, bool) {
	child, okay := res.endpoint.Find(name)
	res.endpoint = child
	return res, okay
}

// Query the endpoint, but do it inside a read lock.
// And, add our turn-metadata to every document
func (res *SessionResource) Query() resource.Document {
	defer res.session.RUnlock()
	res.session.RLock()
	doc := res.endpoint.Query()
	resource.NewDocumentBuilder(&doc).SetMeta(frameKey, res.session.FrameCount())
	return doc
}

// Post to the endpoint, but do it inside a write lock.
func (res *SessionResource) Post(reader io.Reader) (ret resource.Document, err error) {
	// FIX? because the interface takes a reader,
	// we have to sit on the write lock for the duration of the read.
	defer res.session.Unlock()
	res.session.Lock()
	if d, e := res.endpoint.Post(reader); e != nil {
		err = e
	} else {
		resource.NewDocumentBuilder(&d).SetMeta(frameKey, res.session.FrameCount())
		ret = d
	}
	return
}
