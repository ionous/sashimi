package minicon

import (
	"strings"
)

//
// A queue of pending words to display on the MiniCon
// WARNING: collapses multiple spaces and embedded newlines into single spaces;
// and, the hardlines expects that behavior.
//
type PendingWords struct {
	pending []string
}

//
// Forget all pending input.
//
func (this *PendingWords) clear() {
	this.pending = nil
}

//
// Add the passed string, separating the string into words using strings.Fields()
//
func (this *PendingWords) addString(s string) *PendingWords {
	words := strings.Fields(s)
	this.pending = append(this.pending, words...)
	return this
}

//
// Add an explicit new line
//
func (this *PendingWords) newLine() {
	this.pending = append(this.pending, string(LineRune))
}

//
// Extract the oldest word from the queue. Returns the word, and true if there was a valid word.
// NOTE: Explicit new lines are represented as the word: string(LineRune).
//
func (this *PendingWords) popWord() (word string, okay bool) {
	if len(this.pending) > 0 {
		word, okay = this.pending[0], true
		this.pending = this.pending[1:]
	}
	return
}
