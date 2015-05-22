package compiler

import "fmt"

type MultiError struct {
	prev error
	my   error
}

func (e MultiError) Error() string {
	return e.prev.Error() + `\n` + e.my.Error()
}
func AppendError(err error, e error) (ret error) {
	if err == nil {
		ret = e
	} else {
		if e != nil {
			ret = MultiError{err, e}
		} else {
			ret = err
		}
	}
	return ret
}
func PrefixError(err error, s string) (ret error) {
	return fmt.Errorf("%s:%s", s, err)
}

// type ErrorChannel chan error

// func (this ErrorChannel) addError(err error) (hadError bool) {
// 	hadError = err != nil
// 	if hadError {
// 		this <- err
// 	}
// 	return hadError
// }
