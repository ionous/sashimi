package parser

import (
	"fmt"
	"regexp"
	"strings"
)

//
// NormalizeInput turns the passed input to lower-case and replaces multiple spaces with single spaces.
// This helps the parser to match multi-word nouns: "evil   FIsh".
//
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
// SliceArticles returns the noun and its article
//
func SliceArticles(phrase string) (noun string, article string) {
	match := articlesExp.FindStringSubmatch(phrase)
	return match[2], strings.TrimSpace(match[1])
}
