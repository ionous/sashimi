package internal

import (
	S "github.com/ionous/sashimi/source"
	"github.com/ionous/sashimi/util/ident"
)

//
type PropertyRule struct {
	propId ident.Id
	S.PropertyExpectation
}

//
type PendingRules []PropertyRule
