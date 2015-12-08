package facts

import (
	. "github.com/ionous/sashimi/script"
)

func Describe_Facts(s *Script) {
	s.The("kinds", Called("facts"),
		// FIX: interestingly, kinds should have names
		// having the same property as a parent class probably shouldnt be an error
		Have("summary", "text"))

	// FIX: many-to-many doesnt exist; traversing a manually created table of all actors and facts would be fairly heavy; so just using a flag.
	s.The("facts", AreEither("recollected").Or("not recollected").Usually("not recollected"))
}
