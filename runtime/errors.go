package runtime

import "fmt"

//
//
//
type TypeMismatch struct {
	name string
	kind string
}

func (this TypeMismatch) Error() string {
	return fmt.Sprintf("type mismatch %s %s", this.name, this.kind)
}

//
//
//
type NoSuchChoice struct {
	owner  string
	choice string
}

func (this NoSuchChoice) Error() string {
	return fmt.Sprintf("no such choice '%s'.'%s'", this.owner, this.choice)
}

//
// for unexpected runtime errors
//
type RuntimeError struct {
	err error
}

func (this RuntimeError) Error() string {
	return this.err.Error()
}
