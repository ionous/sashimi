package internal

import (
	S "github.com/ionous/sashimi/source"
	"github.com/ionous/sashimi/util/ident"
)

//
type PropertyRule struct {
	fieldName ident.Id
	S.PropertyExpectation
}

//
type PendingRules []PropertyRule
