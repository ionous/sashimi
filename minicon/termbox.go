package minicon

import (
	"fmt"
	"github.com/nsf/termbox-go" // termbox
)

//
// Helper to make termbox a little more object oriented.
//
type TermBox struct {
	enabled bool        // helper for testing
	draw    Pen         // text color
	clear   Pen         // full screen clear color
	terms   TermChannel // non-blocking termbox event management
}

//
// Current foreground and background colors
//
type Pen struct {
	fg, bg termbox.Attribute
}

var defaultPen = Pen{termbox.ColorBlack, termbox.ColorWhite}

//
// After creating a new term box object, the caller must call Init().
//
func NewTermBox() (this *TermBox) {
	return &TermBox{false, defaultPen, defaultPen, make(TermChannel)}
}

//
// Initialize the display. Users should eventually call Close() to return control of the screen to the OS.
//
func (this *TermBox) Init() *TermBox {
	if !this.enabled {
		if e := termbox.Init(); e == nil {
			go this.terms.SendEvents(this)
			this.enabled = true
		}
	}
	return this
}

//
// Kill the terminal are return control of the screen to the OS.
//
func (this *TermBox) Close() {
	if this.enabled {
		termbox.Interrupt() // kill the SendEvents() go routine
		termbox.Close()     // reset the global termbox
	}
}

//
// Clear the screen, hide the cursor.
//
func (this *TermBox) Clear() *TermBox {
	if this.enabled {
		termbox.HideCursor()
		termbox.Clear(this.clear.fg, this.clear.bg)
	}
	return this
}

//
// Blocking waiting for a new terminal event.
//
func (this *TermBox) PollEvent() (evt termbox.Event) {
	if this.enabled {
		evt = termbox.PollEvent()
	} else {
		evt = termbox.Event{Err: fmt.Errorf("not initialized")}
	}
	return evt
}

//
// Check for the next event without blocking.
//
func (this *TermBox) CheckEvent() (ret termbox.Event, okay bool) {
	select {
	// got a new event
	case evt, ok := <-this.terms:
		if !ok {
			ret = termbox.Event{Type: termbox.EventError, Err: fmt.Errorf("channel closed")}
		} else {
			ret = evt
		}
		okay = true // even on error, we have an event
	// nothing pending
	default:
		break
	}
	return ret, okay
}

//
// Return the width and height of the terminal screen.
//
func (this *TermBox) Size() (width int, height int) {
	if this.enabled {
		width, height = termbox.Size()
	} else {
		width, height = 80, 25
	}
	return width, height
}

//
// Set the current draw pen, return the previous draw pen.
//
func (this *TermBox) SetPen(pen Pen) (ret Pen) {
	ret, this.draw = this.draw, pen
	return ret
}

//
// Write a rune at the passed x,y coordinates using the current draw pen
//
func (this *TermBox) SetCell(x, y int, ch rune) {
	if this.enabled {
		termbox.SetCell(x, y, ch, this.draw.fg, this.draw.bg)
	}
}

//
// Place the termbox cursor at the passed x,y coordinates.
//
func (this *TermBox) SetCursor(x, y int) {
	if this.enabled {
		termbox.SetCursor(x, y)
	}
}

//
// Flush any SetCell() changes to the screen so the user can see those changes.
//
func (this *TermBox) Flush() {
	if this.enabled {
		termbox.Flush()
	}
}
