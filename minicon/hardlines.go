package minicon

import "bytes"

//
// Track displayed words and any explicit newlines.
// The interface works character by character to support teletyping.
//
type Hardlines struct {
	pendingWord  bytes.Buffer
	currentLine  []string
	linesOfWords [][]string
}

//
// Add a new word to the current line. Flushes any pending word.
//
func (this *Hardlines) AppendWord(w string) {
	this.StartNewWord() // flush any existing word
	this.currentLine = append(this.currentLine, w)
}

//
// Append a rune to the current word.
//
func (this *Hardlines) AppendRune(ch rune) {
	this.pendingWord.WriteRune(ch)
}

//
// Begin a new word, flushes any pending word.
//
func (this *Hardlines) StartNewWord() {
	if this.pendingWord.Len() > 0 {
		this.currentLine = append(this.currentLine, this.pendingWord.String())
		this.pendingWord.Reset()
	}
}

//
// Terminate the current line, and start a new one. Flushes any pending word.
//
func (this *Hardlines) StartNewLine() {
	this.StartNewWord()
	// flush the current line to the list of all words
	this.linesOfWords = append(this.linesOfWords, this.currentLine)
	this.currentLine = nil
}

//
// `true` if there are no pending words, and nothing has been written to the current line.
//
func (this *Hardlines) IsCurrentLineEmpy() bool {
	return len(this.currentLine) == 0 && this.pendingWord.Len() == 0
}

//
// Total number of explict lines, including the current line ( if non empty ).
//
func (this *Hardlines) NumLines() int {
	cnt := len(this.linesOfWords)
	if !this.IsCurrentLineEmpy() {
		cnt++
	}
	return cnt
}

//
// Return lines with a max length of width, broken at explict new lines favoring word boundries.
//
func (this Hardlines) Reflow(width int) []string {
	output := NewCursorOutput(width)
	if this.NumLines() > 0 {
		lines, currentLine := this.linesOfWords, this.currentLine
		// write full lines
		for _, words := range lines {
			output.AddWords(words)
			output.AddLine()
		}
		// flush the pending word to the current line
		if this.pendingWord.Len() > 0 {
			currentLine = append(currentLine, this.pendingWord.String())
		}
		// write the current line, if any
		if len(currentLine) > 0 {
			output.AddWords(currentLine)
			output.AddSpace()
		}
	}
	return output.Flush()
}
