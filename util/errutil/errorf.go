package errutil

type Errorf interface {
	Errorf(format string, a ...interface{}) error
}
