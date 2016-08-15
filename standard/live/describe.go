package live

import G "github.com/ionous/sashimi/game"

type DescribePhrase struct {
	object string
}

func Describe(object string) DescribePhrase {
	return DescribePhrase{object}
}

func DescribeThe(object G.IObject) DescribePhrase {
	return DescribePhrase{string(object.Id())}
}

func (d DescribePhrase) Execute(g G.Play) {
	if obj := g.The(d.object); obj.Exists() && !obj.Is("scenery") {
		desc := ""
		if obj.Is("unhandled") {
			desc = obj.Text("brief")
		}
		if desc != "" {
			g.Go(Say(desc))
		} else {
			obj.Go("print name")
		}
		obj.Go("print contents")
	}
}
