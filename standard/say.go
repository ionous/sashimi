package standard

import (
	G "github.com/ionous/sashimi/game"
	"strings"
)

const NewLine = "\n"
const OneSpace = " "

func Say(what ...string) SayPhrase {
	return SayPhrase{what: what, sep: NewLine}
}

func (p SayPhrase) OnOneLine() G.RuntimePhrase {
	p.sep = OneSpace
	return p
}

func (p SayPhrase) Execute(g G.Play) {
	text := strings.Join(p.what, p.sep)
	g.Say(text)
}

type SayPhrase struct {
	what []string
	sep  string
}
