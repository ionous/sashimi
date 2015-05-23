package web

import (
	C "github.com/ionous/sashimi/console"
	//"os"
)

// implements IOutput
type SessionOutput struct {
	C.BufferedOutput // implements Print() and Println()
}

func (this *SessionOutput) Write(p []byte) (n int, err error) {
	return 0, err //os.Stderr.Write(p)
}
