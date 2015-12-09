package standard

import (
	G "github.com/ionous/sashimi/game"
	"strings"
)

func Speak(actor string) SpeakPhrase {
	return SpeakPhrase(actor)
}

func (p SpeakPhrase) Says(what ...string) SpeakingPhrase {
	return SpeakingPhrase{actor: string(p), what: what, sep: NewLine}
}

func (p SpeakingPhrase) OnOneLine() G.RuntimePhrase {
	p.sep = OneSpace
	return p
}

func (p SpeakingPhrase) Execute(g G.Play) {
	text := strings.Join(p.what, p.sep)
	g.The(p.actor).Says(text)
}

type SpeakPhrase string

type SpeakingPhrase struct {
	actor string
	what  []string
	sep   string
}
