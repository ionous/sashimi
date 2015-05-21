package minicon

import (
	"time"
)

//
// Preset delays between characters displayed by the Teletype.
//
type PrintSpeed int

const (
	NormalPrinting PrintSpeed = iota
	QuickPrinting
	SkipPrinting
	numPrintSpeeds
)

//
// More faster.
//
func (this PrintSpeed) NextSpeed() (ret PrintSpeed) {
	if ps := this + 1; ps < numPrintSpeeds {
		ret = ps
	} else {
		ret = this
	}
	return ret
}

//
// Duration between letters that are printed by the teletype.
//
func (this PrintSpeed) Duration() time.Duration {
	var printSpeeds = []time.Duration{
		time.Millisecond * 15,
		time.Millisecond / 2,
		0,
	}
	return printSpeeds[this]
}

//
// Block for the Duration().
//
func (this PrintSpeed) Sleep() {
	dur := this.Duration()
	time.Sleep(dur)
}
