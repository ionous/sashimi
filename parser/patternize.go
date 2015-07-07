package parser

import (
	"fmt"
	"regexp"
	"strings"
)

//
// mash the groups together while assigning tags to nouns
//
func newNounCheck(groups []string, tags map[int]bool) (*NounCheck, error) {
	n := &NounCheck{}
	exp, err := n.patternize(groups, tags)
	n.exp = exp
	return n, err
}

func (nouns *NounCheck) patternize(groups []string, tags map[int]bool) (*regexp.Regexp, error) {
	var buffer []string
	for i, content := range groups {
		if !tags[i] {
			// escape any regex characters while keeping bars as "or"
			bars := strings.Split(content, "|")
			if len(bars) <= 1 {
				content = regexp.QuoteMeta(content)
			} else {
				for b, bar := range bars {
					bars[b] = regexp.QuoteMeta(bar)
				}
				s := strings.Join(bars, "|")
				content = "(?:" + s + ")" // non-matching group
			}
			buffer = append(buffer, content)
		} else {
			// match the tag content vs. "something/else"
			tagparts := contentexp.FindStringSubmatch(content)

			// matching failed: so neither "something" nor "something else"
			if len(tagparts) == 0 {
				e := fmt.Errorf("only the tags `something` or `something else` are supported, got `%s` instead.", content)
				return nil, e
			} else {
				// inside the regexp we are building,
				// the only capturing groups are the nouns;
				// the number of nouns we've processed matches the capturing group we're building
				captureIndex := nouns.count + 1

				// the tagparts consists of two strings: the fully matched string, and the matched "else" group.
				if elsePatternIsEmpty := tagparts[1] == ""; elsePatternIsEmpty {
					if e := nouns.addFirstNoun(captureIndex); e != nil {
						return nil, e
					}
				} else {
					if e := nouns.addSecondNoun(captureIndex); e != nil {
						return nil, e
					}
				}
			}
			// replace the tag with a regexp group
			// we preceed with optional articles for nouns
			// leaving them in the matched text so the client code can see the user's intension
			s := fmt.Sprintf("((?:%s |).*?)", articleBar)
			buffer = append(buffer, s)
		}
	}
	// note: see Parser.Parse() which now collapses multiple spaces
	// text := strings.Join(buffer, `\s+`)
	// exp := fmt.Sprintf(`^\s*%s\s*$`, text)
	text := strings.Join(buffer, " ")
	exp := "^" + text + "$"
	return regexp.Compile(exp)
}

// for tokenize:
var contentexp *regexp.Regexp = regexp.MustCompile(`^something($| else)`)
