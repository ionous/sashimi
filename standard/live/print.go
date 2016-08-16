package live

import (
	"fmt"
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/util/lang"
	"strings"
)

// when is the right time for functions versus callbacks?
func ListContents(g G.Play, header string, obj G.IObject) (printed bool) {
	// if something described which is not scenery is on the noun and something which is not the player is on the noun:
	// obviously a filterd callback, visitor, would be nice FilterList("contents", func() ... )
	// FIX: if something has scenery objets, they appear as contents,
	// but then the list is empty. ( ex. lab coat, but it might happen elsewhere )
	// we'd maybe need to know if something printed?
	if contents := obj.ObjectList("contents"); len(contents) > 0 {
		g.Say(header, obj.Text("Name"), "is:")
		for _, content := range contents {
			content.Go("print description")
		}
		printed = true
	}
	return printed
}

func NameFullStop(G.IObject) string {
	return ""
}

type NameStatus func(obj G.IObject) string

func ArticleName(g G.Play, which string, status NameStatus) string {
	return articleName(g, which, false, status)
}
func DefiniteName(g G.Play, which string, status NameStatus) string {
	return articleName(g, which, true, status)
}
func articleName(g G.Play, which string, definite bool, status NameStatus) string {
	obj := g.The(which)
	text := obj.Text("Name")
	if obj.Is("proper-named") {
		text = lang.Titleize(text)
	} else {
		article := ""
		if definite {
			article = "the"
		} else {
			article = obj.Text("indefinite article")
			if len(article) > 0 {
				if obj.Is("plural-named") {
					article = "some"
				} else if lang.StartsWithVowel(text) {
					article = "an"
				} else {
					article = "a"
				}
			}
		}
		text = strings.Join([]string{article, strings.ToLower(text)}, " ")
	}
	if status != nil {
		if s := status(obj); len(s) > 0 {
			text = fmt.Sprintf("%s (%s)", text, s)
		} else {
			text = text + "."
		}
	}
	return text
}
