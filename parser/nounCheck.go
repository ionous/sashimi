package parser

import "fmt"

//
// Helper for reading nouns from patterns.
// ( Patterns are specified by scripters to direct interpretation of player input. )
//
type NounCheck struct {
	firstNoun, secondNoun int
	count, maxNouns       int
}

//
// Create a new noun pattern helper.
//
func newNounCheck(nouns int) *NounCheck {
	return &NounCheck{maxNouns: nouns}
}

//
// Return a lookup table of regexp index => noun index.
//
func (this *NounCheck) matchIndices() []int {
	return []int{this.firstNoun, this.secondNoun}
}

//
// Have all nouns from the pattern been successfully processed?
//
func (this *NounCheck) foundAllNouns() bool {
	return this.count == this.maxNouns
}

//
// The pattern has encountered the first noun.
// Group index: the regexp group index of the current noun.
//
func (this *NounCheck) addFirstNoun(groupIndex int) (err error) {
	if e := this._inc(); e != nil {
		err = e
	} else if this.firstNoun != 0 {
		err = fmt.Errorf("encountered two first nouns")
	} else if this.secondNoun == groupIndex {
		err = fmt.Errorf("duplicate group indexes")
	} else {
		this.firstNoun = groupIndex
	}
	return err
}

//
// The pattern has encountered the second noun.
/// Group index: the regexp group index of the current noun.
//
func (this *NounCheck) addSecondNoun(groupIndex int) (err error) {
	if e := this._inc(); e != nil {
		err = e
	} else if this.secondNoun != 0 {
		err = fmt.Errorf("encountered two second nouns")
	} else if this.firstNoun == groupIndex {
		err = fmt.Errorf("duplicate group indexes")
	} else {
		this.secondNoun = groupIndex
	}
	return err
}

//
// The pattern has encountered a noun.
//
func (this *NounCheck) _inc() (err error) {
	if this.count < this.maxNouns {
		this.count++
	} else {
		err = fmt.Errorf(nounErrors[this.maxNouns])
	}
	return err
}

//
// Error strings for mismatched numbers of nouns.
//
var nounErrors []string = []string{
	"the command wasn't expecting any nouns",
	"the command expected one noun",
	"the command expected two nouns",
	"the command expected three nouns",
}
