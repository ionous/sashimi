package minicon

//
// Status text
//
type Status struct {
	Left, Right string
	Pen         Pen
}

//
// re/draw the status line assuming a window of the passed width.
//
func (this *Status) drawStatus(box *TermBox, width int) {
	oldPen := box.SetPen(this.Pen)
	defer box.SetPen(oldPen)
	// clip the text to the desired width, adding an extra one space margin.
	left, right := " "+this.Left, this.Right+" "
	if lwidth := len(this.Left); lwidth > width {
		left = left[0:width]
	}
	if rwidth := len(this.Right); rwidth > width {
		right = right[0:width]
	}
	y := 0
	// print left
	for i, ch := range left {
		box.SetCell(i, y, ch)
	}
	// when space is limited: the left side overwrites the right
	x := len(left) + 1
	start := width - len(right)
	overlap := x - start
	if overlap > 0 {
		start = x
		right = right[overlap:]
	}
	// print middle:
	for i := x - 1; i < start; i++ {
		box.SetCell(i, y, ' ')
	}
	// print right
	for i, ch := range right {
		box.SetCell(start+i, y, ch)
	}
}
