package parser

import (
	"fmt"
	"regexp"
)

//
// Helper for reading nouns from patterns.
// ( Patterns are specified by scripters to direct interpretation of player input. )
//
type NounCheck struct {
	first, second, count int
	exp                  *regexp.Regexp
}

//
// Return a lookup table of regexp index => noun index.
//
func (nouns *NounCheck) matchIndices() []int {
	return []int{nouns.first, nouns.second}[0:nouns.count]
}

//
// The pattern has encountered the first noun.
// Group index: the regexp group index of the current noun.
//
func (nouns *NounCheck) addFirstNoun(groupIndex int) (err error) {
	if nouns.first != 0 {
		err = fmt.Errorf("encountered two first nouns")
	} else if nouns.second == groupIndex {
		err = fmt.Errorf("duplicate group indexes")
	} else {
		nouns.first = groupIndex
		nouns.count++
	}
	return err
}

//
// The pattern has encountered the second noun.
/// Group index: the regexp group index of the current noun.
//
func (nouns *NounCheck) addSecondNoun(groupIndex int) (err error) {
	if nouns.second != 0 {
		err = fmt.Errorf("encountered two second nouns")
	} else if nouns.first == groupIndex {
		err = fmt.Errorf("duplicate group indexes")
	} else {
		nouns.second = groupIndex
		nouns.count++
	}
	return err
}
