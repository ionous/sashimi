package support

import (
	"log"
	"os/exec"
	"runtime"
	"time"
)

// From Miki Tebeka: Open calls the OS default program for uri
// http://go-wise.blogspot.com/2012/04/open-fileurls-with-default-program.html
func OpenBrowser(uri string) {
	var commands = map[string]string{
		"windows": "start",
		"darwin":  "open",
		"linux":   "xdg-open",
	}
	time.Sleep(300 * time.Millisecond)
	if run, ok := commands[runtime.GOOS]; !ok {
		log.Printf("don't know how to open the browser on %s platform", runtime.GOOS)
	} else {
		cmd := exec.Command(run, uri)
		if e := cmd.Start(); e != nil {
			log.Printf("error opening browser %s", e.Error())
		}
	}
}
