package ess

import (
	"github.com/ionous/sashimi/net/resource"
	"io"
)

const (
	frameKey = "frame"
)

// SessionResource shadows all resources under a sub-tree, locking on all reads and writes, and appending the frame number to returned document's meta-data.
type SessionResource struct {
	session  Session
	endpoint resource.IResource
}

// SessionCreationEndpoint creates new sessions on Post.
func SessionCreationEndpoint(sessions SessionFactory) resource.IResource {
	return &resource.Wrapper{
		Posts: func(_ io.Reader, doc resource.DocumentBuilder) (err error) {
			if session, e := sessions.NewSession(doc); e != nil {
				err = e
			} else {
				doc.SetMeta(frameKey, session.Frame())
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
	resource.NewDocumentBuilder(&doc).SetMeta(frameKey, res.session.Frame())
	return doc
}

// Post to the endpoint, but do it inside a write lock.
// NOTE: the only end point that implements post is likely to be the session end point itself.
func (res *SessionResource) Post(reader io.Reader) (ret resource.Document, err error) {
	// FIX? because the interface takes a reader,
	// we have to sit on the write lock for the duration of the read.
	defer res.session.Unlock()
	res.session.Lock()
	if d, e := res.endpoint.Post(reader); e != nil {
		err = e
	} else {
		resource.NewDocumentBuilder(&d).SetMeta(frameKey, res.session.Frame())
		ret = d
	}
	return
}
