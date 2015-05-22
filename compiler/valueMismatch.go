package compiler

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
)

type ValueMismatch struct {
	id  M.StringId
	was interface{}
	now interface{}
}

func (this ValueMismatch) Error() string {
	return fmt.Sprintf("value mismatch: %s had %v requested %v", this.id, this.was, this.now)
}
