package errutil

// Append joins two errors into one; either or both can be nil.
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

// multiError chains errors together.
// it can expand "recursively", allowing an chain of errors.
type multiError struct {
	prev error
	my   error
}

// Error returns the combination of errors separated by a newline.
func (e multiError) Error() string {
	return e.prev.Error() + "\n" + e.my.Error()
}
