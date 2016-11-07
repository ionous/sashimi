package lang

import (
	"bitbucket.org/pkg/inflect"
	"regexp"
	"strings"
)

var Articles = []string{"the", "a", "an", "our", "some"}
var articleBar = strings.Join(Articles, "|")
var articles = regexp.MustCompile(`^((?i)` + articleBar + `)\s`)
var articleBare = regexp.MustCompile("^(" + articleBar + ")$")

const NewLine = "\n"
const Space = " "

func IsArticle(s string) bool {
	return articleBare.MatchString(s)
}

func SliceArticle(str string) (article, bare string) {
	n := strings.TrimSpace(str)
	if pair := articles.FindStringIndex(n); pair == nil {
		bare = n
	} else {
		split := pair[1] - 1
		article = n[:split]
		bare = strings.TrimSpace(n[split:])
	}
	return article, bare
}
func StripArticle(str string) string {
	_, bare := SliceArticle(str)
	return bare
}

//
var Singularize = inflect.Singularize

//
var Pluralize = inflect.Pluralize

// Capitalize returns a new string, starting the first word with a capital.
var Capitalize = inflect.Capitalize

// Titleize returns a new string, starting every word with a capital.
var Titleize = inflect.Titleize

func StartsWith(s string, set ...string) (ok bool) {
	for _, x := range set {
		if strings.HasPrefix(s, x) {
			ok = true
			break
		}
	}
	return ok
}

//http://www.mudconnect.com/SMF/index.php?topic=74725.0
func StartsWithVowel(str string) (vowelSound bool) {
	s := strings.ToUpper(str)
	if StartsWith(s, "A", "E", "I", "O", "U") {
		if !StartsWith(s, "EU", "EW", "ONCE", "ONE", "OUI", "UBI", "UGAND", "UKRAIN", "UKULELE", "ULYSS", "UNA", "UNESCO", "UNI", "UNUM", "URA", "URE", "URI", "URO", "URU", "USA", "USE", "USI", "USU", "UTA", "UTE", "UTI", "UTO") {
			vowelSound = true
		}
	} else if StartsWith(s, "HEIR", "HERB", "HOMAGE", "HONEST", "HONOR", "HONOUR", "HORS", "HOUR") {
		vowelSound = true
	}
	return vowelSound
}
