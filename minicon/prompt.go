package minicon

import (
	"bytes"
)

//
// User command managment.
//
type Prompt struct {
	buffer bytes.Buffer
	cursor CursorWidth
	y      int
}

//
// Each MiniCon has its own prompt.
//
func NewPrompt() *Prompt {
	return new(Prompt)
}

//
// Clear and return the user's previous input as a string.
//
func (this *Prompt) Clear() string {
	s := this.buffer.String()
	this.buffer.Reset()
	return s
}

//
// Repaint the prompt, input string, and cursor.
//
func (this *Prompt) RefreshPrompt(box *TermBox) *TermBox {
	width, height := box.Size()
	this.y, this.cursor = height-1, CursorWidth{1, 1, width}
	return this.refresh(box)
}

//
// Ensure the termbox's cursor is on the prompt line at the proper point.
//
func (this *Prompt) RefreshCursor(box *TermBox) *TermBox {
	box.SetCursor(this.cursor.pos, this.y)
	return box
}

//
// Set the display input to the passed string.
//
func (this *Prompt) SetInput(box *TermBox, str string) *TermBox {
	this.cursor.ResetCursor()
	this.buffer.Truncate(0)
	this.buffer.WriteString(str)
	return this.refresh(box)
}

//
// Add a rune to the user's input.
//
func (this *Prompt) WriteRune(box *TermBox, ch rune) *TermBox {
	this.buffer.WriteRune(ch)
	this._put(box, ch, 1)
	return this.RefreshCursor(box)
}

//
// Remove the most recent rune from the user's input.
//
func (this *Prompt) DeleteRune(box *TermBox) *TermBox {
	if l := this.buffer.Len(); l > 0 {
		this.buffer.Truncate(l - 1)
		this._put(box, SpaceRune, -1)
		this._put(box, SpaceRune, 0)
	}
	return this.RefreshCursor(box)
}

//
// repaint the prompt, input string, and cursor.
//
func (this *Prompt) refresh(box *TermBox) *TermBox {
	// draw the prompt
	y := this.y
	box.SetCell(0, y, '>')
	// clear the input region
	for x := this.cursor.minPos; x < this.cursor.width; x++ {
		box.SetCell(x, y, SpaceRune)
	}
	// print user string
	str := this.buffer.String()
	if len(str) <= 0 {
		// set the cursor if nothing else
		this._put(box, 0, 0)
	} else {
		for _, ch := range str {
			this._put(box, ch, 1)
		}
	}
	return this.RefreshCursor(box)
}

//
// helper for displaying characters.
//
func (this *Prompt) _put(box *TermBox, ch rune, step int) {
	box.SetCell(this.cursor.pos, this.y, ch)
	if this.cursor.IsStepInLimits(step) {
		this.cursor.StepCursor(step)
	}
}
