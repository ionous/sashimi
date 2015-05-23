package source

import (
	"fmt"
)

// statement options
type Options map[string]string

func (this Options) Error() string {
	return fmt.Sprintf("unknown instance options specified %s", this)
}
