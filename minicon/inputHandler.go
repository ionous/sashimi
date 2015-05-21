package minicon

import (
	"github.com/nsf/termbox-go" // termbox; for event handling
)

//
// Helper to processes input events to various useful ends.
//
type InputHandler struct {
	prompt   *Prompt       // input adds to the user's buffered, pending command
	box      *TermBox      // input displays on the screen
	teletype *Teletype     // input can advance the speed of the teletype
	history  *History      // input can add and cycle through history
	marker   HistoryMarker // restores history if the user cycles through history
}

//
// Every unqiue user entry needs its own unique InputHandler().
//
func NewInputHandler(
	input *Prompt,
	box *TermBox,
	teletype *Teletype,
	history *History,
) *InputHandler {
	return &InputHandler{input, box, teletype, history, history.Mark()}
}

//
// Return history to most recent moment in time. ( In case, via user input, the user cycled through old commands. )
//
func (this *InputHandler) RestoreHistory() {
	this.history.Restore(this.marker)
}

//
// Process a single terminal event:
// read user input, add characters to the prompt, cycle through history.
// Returns the user's command if they typed text and pressed enter; empty string otherwise.
// Returns `true` if the terminal screen needs a refresh ( ex. a resize event was detected. )
//
func (this *InputHandler) HandleTermEvent(evt termbox.Event,
) (userInput string, refresh bool,
) {
	prompt, box := this.prompt, this.box
	switch evt.Type {
	case termbox.EventKey:
		if evt.Ch != 0 {
			prompt.WriteRune(box, evt.Ch).Flush()
		} else {
			switch evt.Key {
			case termbox.KeyDelete,
				termbox.KeyBackspace,
				termbox.KeyBackspace2:
				prompt.DeleteRune(box).Flush()

			// history:
			case termbox.KeyArrowUp:
				// if the items back in time have been exhausted, str is blank.
				// don't clear the input, stick to that last item.
				if str, ok := this.history.Back(); ok {
					prompt.SetInput(box, str).Flush()
				}

			case termbox.KeyArrowDown:
				// if the items forward have been exhausted, str is blank.
				// it's okay, good even, to restore the input to a blank string
				str, _ := this.history.Forward()
				prompt.SetInput(box, str).Flush()

			// got valid input
			case termbox.KeyEnter:
				userInput = prompt.Clear()
				if userInput != "" {
					this.marker = this.history.Add(userInput, this.marker)
				}

			case termbox.KeySpace:
				prompt.WriteRune(box, SpaceRune).Flush()

			case termbox.KeyEsc, termbox.KeyTab:
				this.teletype.NextSpeed()
			}
		}
	case termbox.EventResize:
		refresh = true

		//case termbox.EventError, termbox.EventInterrupt:
		//	break WaitLoop
		//case termbox.EventMouse:
		//case termbox.EventRaw:
	}
	return userInput, refresh
}
