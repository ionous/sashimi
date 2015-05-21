package minicon

//
// Tracks the screen position of where text will next output.
//
type CursorWidth struct {
	pos, minPos, width int
}

//
// Is the cursor at the start of its range?
//
func (this CursorWidth) IsStart() bool {
	return this.pos == this.minPos
}

//
// Number of valid positions remaining in the range.
//
func (this CursorWidth) IsValid() bool {
	return this.IsStepInLimits(0)
}

//
// If the cursor were advanced by the passed number of steps,
// would the output position be within the current bounds?
//
func (this CursorWidth) IsStepInLimits(step int) bool {
	next := this.pos + step
	return next >= this.minPos && next < this.width
}

//
// Number of valid positions remaining in the range.
//
func (this CursorWidth) RemainingSteps() int {
	return this.width - this.pos
}

//
// Move the cursor back to the start of its range.
//
func (this *CursorWidth) ResetCursor() {
	this.pos = this.minPos
}

//
// Move the cursor, return true if the new position IsValid()
//
func (this *CursorWidth) StepCursor(step int) bool {
	this.pos += step
	return this.IsValid()
}
