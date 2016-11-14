package script

import (
	. "github.com/ionous/sashimi/script"
)

// FIX: replace  with player
func DirectlyFollows(other string) IFragment {
	return discuss("directly following", other)
}
func IndirectlyFollows(other string) IFragment {
	return discuss("indirectly following", other)
}

func IsPermittedBy(fact string) IFragment {
	return requires("permitted", fact)
}

func IsProhibitedBy(fact string) IFragment {
	return requires("prohibited", fact)
}

func discuss(how, other string) IFragment {
	// FIX: a way to change the orgin?
	return NewFunctionFragment(func(s *S.Statements) error {
		b.The("following quips",
			Table("following", "indirectly following", "leading").Has( //-property
				b.Subject(), how, other))
		return nil
	})
}

func requires(how, fact string) IFragment {
	return NewFunctionFragment(func(s *S.Statements) error {
		b.The("quip requirements",
			Table("fact", "permitted", "quip"). //-property
								Has(fact, how, b.Subject()))
		return nil
	})
}
