package minicon

import (
	"fmt"
	"github.com/nsf/termbox-go" // termbox; for event handling
)

//
// A simple text window with status line and user input prompt.
//
type MiniCon struct {
	box       *TermBox // terminal window
	Status    Status   // status line at the top of the window
	prompt    *Prompt
	pending   PendingWords // printed but not yet displayed
	teletype  Teletype     // displays printed words character by character
	hardlines Hardlines    // displayed words broken into lines
	padding   int          // scroll up, down
	history   *History     // command input history
}

// NOTE: Users should call Close() on the console when they are done
func NewMiniCon() *MiniCon {
	box := NewTermBox().Init()
	statusPen := Pen{termbox.ColorWhite, termbox.ColorBlack}
	disp := &MiniCon{
		box:     box,
		Status:  Status{Pen: statusPen},
		prompt:  NewPrompt(),
		history: NewHistory(25),
	}
	return disp
}

//
// Restore the user's screen and free resources
//
func (this *MiniCon) Close() {
	this.box.Close()
}

//
// Clear all displayed and printed text; but not status and not history.
//
func (this *MiniCon) ClearScreen() {
	this.box.Clear()
	this.pending.clear()
	this.teletype.Clear()
	this.hardlines = Hardlines{}
	this.padding = 0
}

//
// Add the passed arguements as space separated words as pending output.
// NOTE: the output is not displayed until Update()
//
func (this *MiniCon) Print(args ...interface{}) {
	line := fmt.Sprint(args...)
	this.pending.addString(line)
}

//
// Same as Print() with a new line appended.
//
func (this *MiniCon) Println(args ...interface{}) {
	line := fmt.Sprint(args...)
	this.pending.addString(line).newLine()
}

//
// Flush any pending teletype text to the screen
//
func (this *MiniCon) Flush() {
	this.teletype.PrintSpeed = SkipPrinting
	for !this.advanceTeletype() {
	}
	this.RefreshDisplay().Flush()
	this.teletype.PrintSpeed = NormalPrinting
}

const (
	leftMargin, rightMargin, topMargin, bottomMargin = 0, 0, 2, 1
)

//
// Redraw the screen, reflowing all text to accomdate the window's current dimensions.
//
func (this *MiniCon) RefreshDisplay() *TermBox {
	box := this.box
	fullwidth, fullheight := box.Clear().Size()

	this.Status.drawStatus(box, fullwidth)
	width := fullwidth - leftMargin - rightMargin
	height := fullheight - topMargin - bottomMargin

	// re/format the hardlines --
	// we could cache this if that was useful.
	realLines := this.hardlines.Reflow(width)
	totalLines := len(realLines) + this.padding
	lastX, lastY := leftMargin, topMargin-1
	if totalLines > 0 {
		// how many of the desired lines fit on screen?
		visibleLines := height
		if totalLines < visibleLines {
			visibleLines = totalLines
		}
		// subsection the real lines to just what we want to see...
		topLine := totalLines - visibleLines
		if topLine >= 0 {
			lines := realLines[topLine:]
			for y, line := range lines {
				y := y + topMargin
				for x, ch := range line {
					x := x + leftMargin
					this.box.SetCell(x, y, ch)
					lastX = x
				}
				lastY = y
			}
		}
	}
	// calc the position for the teletype
	var x, y int
	if !this.hardlines.IsCurrentLineEmpy() {
		x = lastX + 1 // space
		y = lastY
	} else {
		x = leftMargin
		y = lastY + 1
	}
	this.teletype.Resize(x, y, leftMargin, topMargin, width, height)
	return this.prompt.RefreshPrompt(box)
}

//
// Display pending output and read input from the user.
// Blocks until the user has entered some text and pressed return.
//
func (this *MiniCon) Update() (ret string) {
	box, prompt := this.box, this.prompt

	// user input buffer and user input handler
	inputHandler := NewInputHandler(prompt, box, &this.teletype, this.history)
	defer inputHandler.RestoreHistory()
	prompt.RefreshPrompt(box).Flush()

WaitLoop:
	for {
		if allDone := this.advanceTeletype(); allDone {
			this.teletype.PrintSpeed = NormalPrinting
			this.RefreshDisplay().Flush()
		}

		if evt, ok := this.box.CheckEvent(); ok {
			userInput, refresh := inputHandler.HandleTermEvent(evt)
			if refresh {
				this.RefreshDisplay().Flush()
			}
			if userInput != "" {
				ret = userInput
				this.teletype.PrintSpeed = NormalPrinting
				break WaitLoop
			}
		}
	}
	return ret
}

//
// update the teletype
//
func (this *MiniCon) advanceTeletype() (allDone bool) {
	box := this.box
	// NOTE: display words sleeps slightly depending on the teletype speed.
	if ch, okay := this.teletype.TypeNextRune(box); okay {
		this.hardlines.AppendRune(ch)
		if this.teletype.PrintSpeed != SkipPrinting {
			box.Flush() // update the screen with the rune we printed
		}
	} else {
		if word, ok := this.pending.popWord(); !ok {
			// done with all pending words?
			allDone = true
		} else {
			// a word consisting only of a newline is a special, explict, request for a new line.
			if word == string(LineRune) {
				this.hardlines.StartNewLine()
				this.teletype.NewLine()
			} else {
				this.hardlines.StartNewWord()
				this.teletype.QueueWord(word)
			}
			// moving to a new line can leave our y-position invalid
			remain := this.teletype.RemainingLines()
			if remain <= 0 {
				this.padding = 1
				// full refresh
				if this.teletype.PrintSpeed != SkipPrinting {
					this.RefreshDisplay().Flush()
				}
			} else if this.padding > remain {
				this.padding = remain
			}
		}
	}
	return allDone
}
