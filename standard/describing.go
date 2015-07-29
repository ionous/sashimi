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
	obj := g.The(d.object)
	if !obj.Is("scenery") {
		desc := ""
		if obj.Is("unhandled") {
			desc = obj.Text("brief")
		}
		if desc != "" {
			g.Say(desc)
		} else {
			obj.Go("print name")
		}
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
		s.The("containers",
			//print description
			When("describing").Always(func(g G.Play) {
				this := g.The("action.Source")
				if (this.Is("open") || this.Is("transparent")) && !this.Is("scenery-content") {
					g.Say(" ")
					g.Go(DescribeThe(this))
					listContents(g, "In the", this)
					g.StopHere()
				}
			}))

		s.The("supporters",
			When("describing").Always(func(g G.Play) {
				this := g.The("supporter")
				g.Say(" ")
				g.Go(DescribeThe(this))
				listContents(g, "On the", this)
				g.StopHere()
			}))
	})
}
