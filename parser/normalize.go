package parser

import (
	"fmt"
	"regexp"
	"strings"
)

// lower case the input,
// strip multiple spaces and replace with single spaces,
// [ this helps match multi-word nouns: "evil   fish" ]
func NormalizeInput(input string) string {
	lower := strings.ToLower(input)
	fields := strings.Fields(lower)
	return strings.Join(fields, " ")
}

// see also patternize()
var articles = []string{"the", "a", "an", "our", "some"}
var articleBar = strings.Join(articles, " |")
var articlesExp = regexp.MustCompile(fmt.Sprintf("^(%s |)(.*?)$", articleBar))

//
// returns noun, article of the noun
//
func NormalizeNoun(noun string) (string, string) {
	match := articlesExp.FindStringSubmatch(noun)
	return match[2], strings.TrimSpace(match[1])
}
