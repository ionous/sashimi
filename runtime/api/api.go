package api

import "github.com/ionous/sashimi/util/ident"

type Model interface {
	GetAction(ident.Id) (Action, bool)
	GetParserActions(func(act Action, commands []string) bool)
}

type Action interface {
	GetId() ident.Id
	GetActionName() string
	GetEventName() string
	GetNouns() Nouns
}

type NounType int

const (
	SourceNoun NounType = iota
	TargetNoun
	ContextNoun
)

// a list of nouns
type Nouns []ident.Id

func (n Nouns) GetNounCount() int {
	return len(n)
}

func (n Nouns) Get(t NounType) (ret ident.Id) {
	if int(t) < len(n) {
		ret = n[t]
	}
	return ret
}
