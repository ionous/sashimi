package resource

import (
	"fmt"
	"strings"
)

// FindResource expands the passed resource, using each element of the passed path in turn.
// Returns an error, PathError, describing the extent of the matched path.
func FindResource(res IResource, path string) (ret IResource, err error) {
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if r, ok := res.Find(part); ok {
			res = r // update res for next loop
		} else {
			err = PathError{parts, i}
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

func (err PathError) Error() string {
	return fmt.Sprintf("failed to find resource %d(%s) in %s", err.FailedIndex,
		err.Parts[err.FailedIndex],
		strings.Join(err.Parts, "/"))
}
