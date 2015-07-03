package runtime

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
)

//
// Return a noun which matches an instance's string id
//
type NounFactory struct {
	act   *M.ActionInfo
	model *M.Model
	rank  int
}

//
// Called successively for each word in the user input.
//
func (this *NounFactory) MatchNoun(word string, _ string) (noun string, err error) {
	this.rank++
	nouns := this.act.NounSlice()
	class := nouns[this.rank]
	if class == nil {
		err = fmt.Errorf("You've told me more than I've understood.")
	} else {
		names := this.model.NounNames[word]
		if len(names) == 0 {
			err = fmt.Errorf("I don't see any such thing.")
		} else {
			for _, name := range names {
				if inst, ok := this.model.Instances[name]; ok {
					if inst.Class().CompatibleWith(class.Id()) {
						noun = name.String()
						break
					}
				}
			}
			if noun == "" {
				err = fmt.Errorf("I don't know how to use that for this.")
			}
		}
	}
	return noun, err
}
