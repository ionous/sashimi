package resource

import (
	"fmt"
	"io"
)

//
// Turn one or more IResource compatible functions into a full interface implementation.
//
type Wrapper struct {
	Finds   func(string) (IResource, bool)
	Queries func(DocumentBuilder)
	Posts   func(io.Reader, DocumentBuilder) error
}

func (w Wrapper) Find(child string) (ret IResource, okay bool) {
	if finds := w.Finds; finds != nil {
		ret, okay = finds(child)
	}
	return
}

func (w Wrapper) Query() (ret Document) {
	if queries := w.Queries; queries != nil {
		queries(DocumentBuilder{&ret})
	}
	return
}

func (w Wrapper) Post(reader io.Reader) (ret Document, err error) {
	if posts := w.Posts; posts != nil {
		err = posts(reader, DocumentBuilder{&ret})
	} else {
		err = fmt.Errorf("not supported")
	}
	return
}
