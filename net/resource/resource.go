package resource

import (
	"io"
)

//
// A rest-ish endpoint.
//
type IResource interface {
	// Return the named sub-resource
	Find(string) (IResource, bool)
	// Return the resource
	Query() Document
	// Post to the resource
	Post(io.Reader) (Document, error)
}
