package standard

import (
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/util/lang"
	"strings"
)

func Say(what ...string) SayPhrase {
	return SayPhrase{what: what, sep: lang.NewLine}
}

func (p SayPhrase) OnOneLine() G.RuntimePhrase {
	p.sep = lang.Space
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
