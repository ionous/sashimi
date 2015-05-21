package minicon

//
// Prints words one character at a time while word wrapping to the current region.
//
type Teletype struct {
	_cursor    CursorRegion
	PrintSpeed PrintSpeed
	pending    []rune
	//cursed     *TermBox
}

func (this *Teletype) Resize(x, y, left, top, width, height int) {
	this._cursor.Reset(x, y, left, top, width, height)
	//this.cursed = nil
}

//
// Abandon any pending characters.
//
func (this *Teletype) Clear() {
	this.pending = nil
	//this.cursed = nil
}

//
// Return the amount of vertical space remaining.
//
func (this *Teletype) RemainingLines() int {
	return this._cursor.y.RemainingSteps()
}

//
// Advance to the next PrintSpeed.
//
func (this *Teletype) NextSpeed() {
	this.PrintSpeed = this.PrintSpeed.NextSpeed()
}

//
// Advance to the next line.
//
func (this *Teletype) NewLine() {
	// if this.cursed != nil {
	// 	cursor := &this._cursor
	// 	this.cursed.SetCell(cursor.x.pos+1, cursor.y.pos, SpaceRune)
	// 	this.cursed = nil
	// }
	this._cursor.NewLine()
}

//
// Queues a new word for printing.
//
func (this *Teletype) QueueWord(word string) {
	// wont fit, so first move to a new line
	if okay := this._cursor.x.IsStepInLimits(len(word)); !okay {
		this.NewLine()
	}
	this.pending = append(this.pending, []rune(word)...)
}

//
// Prints a single rune from the most recently QueuedWord to the display.
// Returns the rune printed, and true if the rune there was a rune to print.
//
func (this *Teletype) TypeNextRune(box *TermBox) (ret rune, okay bool) {
	if ok := len(this.pending) > 0; ok {
		cursor := &this._cursor
		// pop a pending rune
		ch, remaining := this.pending[0], this.pending[1:]
		this.pending = remaining

		// move the cursor, and draw the rune at the cursor's last position
		cursor.Write(box, ch)

		// after every word, leave a space so the cursor can immediately print new chars
		if cursor.x.IsValid() && len(remaining) == 0 && !cursor.x.IsStart() {
			cursor.Write(box, SpaceRune)
		}

		// NOTE: when printing at top speed: ignore ignore the cursor, flush, and sleep.
		if this.PrintSpeed.Duration() > 0 {
			// if cursor.x.IsStepInLimits(1) {
			// 	box.SetCell(cursor.x.pos+1, cursor.y.pos, '_')
			// 	this.cursed = box
			// }
			box.Flush()
			this.PrintSpeed.Sleep()
		}

		// setup the return value
		ret, okay = ch, true
	}
	return ret, okay
}
