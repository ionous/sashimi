package fact

// FIX: replace  with player, go learn
// ALSO: if this were in the "fact" package, it could be: fact.Learn
// and maybe prop.Give?
func Learn(fact string) FactPhrase {
	return FactPhrase{fact}
}
func LearnThe(fact G.IObject) FactPhrase {
	return FactPhrase{fact.Id().String()}
}

type FactPhrase struct {
	fact string
}

func (p FactPhrase) Execute(g G.Play) {
	PlayerMemory(g).Learns(g.The(p.fact))
}

func PlayerLearns(g G.Play, fact string) (newlyLearned bool) {
	return PlayerMemory(g).TriesToLearn(g.The(fact))
}

func Describe_Facts(s *Script) {
	s.The("kinds", Called("facts"),
		// FIX: interestingly, kinds should have names
		// having the same property as a parent class probably shouldnt be an error
		Have("summary", "text"))

	// FIX: should be "data"
	s.The("kinds", Called("quip requirements"),
		Have("fact", "fact"),
		AreEither("permitted").Or("prohibited"),
		Have("quip", "quip"))
}
