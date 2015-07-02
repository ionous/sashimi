package errutil

// Append joins two errors into one; either or both errors can be nil.
func Append(err error, e error) (ret error) {
	if err == nil {
		ret = e
	} else if e != nil {
		ret = multiError{err, e}
	} else {
		ret = err
	}
	return ret
}

// MultiError joins two errors into one.
// It can expand "recursively", allowing an chain of errors.
type multiError struct {
	prev error
	my   error
}

// Error returns a string containing MultiError's two errors, separated by a newline.
func (e multiError) Error() string {
	return e.prev.Error() + "\n" + e.my.Error()
}
