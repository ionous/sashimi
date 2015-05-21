package minicon

import (
	"bytes"
)

// Reflow, hard line break helper.
type CursorOutput struct {
	lines   []string
	pending *bytes.Buffer // current line accumulation
	CursorWidth
}

//
// Return all lines as an array of lines.
//
func NewCursorOutput(width int) CursorOutput {
	return CursorOutput{nil, new(bytes.Buffer), CursorWidth{0, 0, width}}
}

//
// Return all lines as an array of lines.
//
func (this *CursorOutput) Flush() []string {
	lines := this.lines
	line := this.pending.String()
	if len(line) > 0 {
		lines = append(lines, line)
	}
	return lines
}

//
// Finish the pending line, and start a new line.
//
func (this *CursorOutput) AddLine() {
	line := this.pending.String()
	this.lines = append(this.lines, line)
	this.pending = new(bytes.Buffer)
	this.ResetCursor()
}

//
// Inject a space between two words
//
func (this *CursorOutput) AddSpace() {
	this.AddRune(SpaceRune)
}

//
// Write space separated words, breaking on lines as necessary.
//
func (this *CursorOutput) AddWords(words []string) {
	for i, word := range words {
		if i != 0 {
			this.AddSpace()
		}
		this.AddWord(word)
	}
}

// Write a single word, breaking on lines as necessary.
func (this *CursorOutput) AddWord(word string) {
	// wont fit, so first move to a new line
	if !this.IsStepInLimits(len(word)) {
		this.AddLine()
	}
	// and, print the word:
	for _, ch := range word {
		this.AddRune(ch)
	}
}

// Write a char, possibly leaving the cursor at the start of a new, empty line.
func (this *CursorOutput) AddRune(ch rune) {
	this.pending.WriteRune(ch)
	if !this.StepCursor(1) {
		this.AddLine()
	}
}
