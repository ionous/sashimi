package resource

import (
	"io"
)

//
// IResource interfaces with a rest-ish endpoint.
// See also, Wrapper, which provides a function-based adapter.
//
type IResource interface {
	// Return the named sub-resource
	Find(string) (IResource, bool)
	// Return the resource
	Query() Document
	// Post to the resource
	Post(io.Reader) (Document, error)
}
