package standard

import (
	G "github.com/ionous/sashimi/game"
	"strings"
)

const NewLine = "\n"
const OneSpace = " "

func Say(s ...string) SayPhrase {
	return SayPhrase{s: s, sep: NewLine}
}

func (p SayPhrase) OnOneLine() G.RuntimePhrase {
	p.sep = OneSpace
	return p
}

func (p SayPhrase) Execute(g G.Play) {
	text := strings.Join(p.s, p.sep)
	g.Say(text)
}

type SayPhrase struct {
	s   []string
	sep string
}
