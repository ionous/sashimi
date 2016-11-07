package parser

import (
	"github.com/ionous/sashimi/util/lang"
	"strings"
)

// NormalizeInput turns the passed input to lower-case and replaces multiple spaces with single spaces.
// This helps the parser to match multi-word nouns: "evil   FIsh".
func NormalizeInput(input string) string {
	lower := strings.ToLower(input)
	fields := strings.Fields(lower)
	next := []string{}
	for _, f := range fields {
		if !lang.IsArticle(f) {
			next = append(next, f)
		}
	}
	return strings.Join(next, " ")
}
