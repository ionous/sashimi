package console

import (
	"io/ioutil"
	"log"
	"os"
)

//
func NewLogger(verbose bool) (logger *log.Logger) {
	if verbose {
		logger = log.New(os.Stderr, "test:", log.Lshortfile)
	} else {
		logger = log.New(ioutil.Discard, "", 0)
	}
	return logger
}
