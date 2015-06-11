package resource

import (
	"fmt"
	"strings"
)

//
func NewResourcePath(path string) (ret ResourcePath) {
	return ResourcePath{strings.Split(path, "/")}
}

// helper to find a resource.
type ResourcePath struct {
	parts []string
}

//
// Find a resource endpoint, using the passed resource as a starting point and matching all the elements of this path.
// Returns an error describing the extent of the matched path
//
func (this ResourcePath) FindResource(res IResource) (ret IResource, err error) {
	for i, part := range this.parts {
		if r, ok := res.Find(part); ok {
			res = r // update res for next loop
		} else {
			err = PathError{this.parts, i}
			break
		}
	}
	if err == nil {
		ret = res
	}
	return res, err
}

type PathError struct {
	Parts       []string
	FailedIndex int
}

func (this PathError) Error() string {
	return fmt.Sprintf("failed to find resource %d(%s) in %s", this.FailedIndex,
		this.Parts[this.FailedIndex],
		strings.Join(this.Parts, "/"))
}
