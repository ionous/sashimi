package minicon

//
// Helper to write left-to-right, bottom-to-top within a specified bounds.
//
type CursorRegion struct {
	x, y CursorWidth
}

// Define a new valid region for a cursor passing starting x, starting y, and the region bounds.
func (this *CursorRegion) Reset(x, y, left, top, width, height int) {
	this.x = CursorWidth{x, left, width}
	this.y = CursorWidth{y, top, height}
}

// Move to the start of the next line.
func (this *CursorRegion) NewLine() {
	if this.y.IsValid() {
		this.x.ResetCursor()
		this.y.StepCursor(1)
	}
}

// Advance the space of one rune.
// Returns the previous position of the cursor, and 'true' if the cursor was left in a valid state.
func (this *CursorRegion) Write(box *TermBox, ch rune) bool {
	okay := this.x.IsValid() && this.y.IsValid()
	if okay {
		box.SetCell(this.x.pos, this.y.pos, ch)
		this.x.StepCursor(1)
	}
	return okay
}
