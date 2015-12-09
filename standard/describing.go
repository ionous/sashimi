package standard

import (
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
)

type DescribePhrase struct {
	object string
}

func Describe(object string) DescribePhrase {
	return DescribePhrase{object}
}

func DescribeThe(object G.IObject) DescribePhrase {
	return DescribePhrase{object.Id().String()}
}

func (d DescribePhrase) Execute(g G.Play) {
	if obj := g.The(d.object); obj.Exists() && !obj.Is("scenery") {
		desc := ""
		if obj.Is("unhandled") {
			desc = obj.Text("brief")
		}
		var ch G.IPromise
		if desc != "" {
			ch = g.Go(Say(desc))
		} else {
			ch = obj.Go("print name")
		}
		ch.Then(func(g G.Play) {
			obj.Go("print contents")
		})
	}
}

func init() {
	AddScript(func(s *Script) {
		s.The("objects",
			Can("print description").And("describing").RequiresNothing(),
			To("print description", func(g G.Play) {
				g.Go(Describe("object"))
			}))

		// FIX: When() puts the contents after the object
		// look at some default actions of the DOM
		// maybe it'd be better to put the print, not in the action,
		// but in a target handler: then this could be after by being in the capture.

		// FIX: After() isnt working well, it goes into the default action
		// but not all objects are containers, so it errors
	})
}
