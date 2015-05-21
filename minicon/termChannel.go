package minicon

import (
	"github.com/nsf/termbox-go" // termbox
)

//
// Helper for communicating with termbox in a non blocking way.
//
type TermChannel chan termbox.Event

//
// Pushes all events, one at a time, into the passed channel.
// Closes the channel and returns only if termbox.Interrupt() was called.
//
func (this TermChannel) SendEvents(box *TermBox) {
	for {
		evt := box.PollEvent()
		if evt.Type == termbox.EventInterrupt {
			close(this)
			break
		}
		this <- evt
	}
}
