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

func (this Wrapper) Find(child string) (ret IResource, okay bool) {
	if finds := this.Finds; finds != nil {
		ret, okay = finds(child)
	}
	return
}

func (this Wrapper) Query() (ret Document) {
	if queries := this.Queries; queries != nil {
		queries(DocumentBuilder{&ret})
	}
	return
}

func (this Wrapper) Post(reader io.Reader) (ret Document, err error) {
	if posts := this.Posts; posts != nil {
		err = posts(reader, DocumentBuilder{&ret})
	} else {
		err = fmt.Errorf("not supported")
	}
	return
}
