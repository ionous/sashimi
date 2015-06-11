package session

import "fmt"

// Error signal that the session has finished for some reason
type SessionClosed struct {
	Reason string
}

func (this SessionClosed) Error() string {
	return fmt.Sprintf("Session Closed: %s", this.Reason)
}
