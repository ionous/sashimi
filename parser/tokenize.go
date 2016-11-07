package parser

import (
	"regexp"
	"strings"
)

var tokenexp *regexp.Regexp = regexp.MustCompile(`({{.*?}})`)

// split the passed string into groups separated by {{bracketed text}}
// tags provide the indicies of those text groups
func tokenize(text string) (groups []string, tags map[int]bool) {
	if len(text) > 0 {
		tags = make(map[int]bool)
		// find all pairs of brackets:
		tagIndexes, textStart := tokenexp.FindAllStringIndex(text, -1), 0
		for _, tagRange := range tagIndexes {
			tagStart, tagEnd := tagRange[0], tagRange[1]
			// extract the preceeding fields
			prev := text[textStart:tagStart]
			fields := strings.Fields(prev)
			groups = append(groups, fields...)

			// add the contents of the tag; removing duplicate spaces inside the tag
			tags[len(groups)] = true
			content := text[tagStart+2 : tagEnd-2]
			tag := strings.Join(strings.Fields(content), " ")
			groups = append(groups, tag)
			// advance the position in the input string from which we're reading
			textStart = tagEnd
		}
		// add any text after the last tag
		fini := strings.TrimSpace(text[textStart:])
		if len(fini) > 0 {
			groups = append(groups, fini)
		}
	}
	return groups, tags
}
