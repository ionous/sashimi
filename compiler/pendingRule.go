package compiler

import (
	M "github.com/ionous/sashimi/model"
	S "github.com/ionous/sashimi/source"
)

//
type PropertyRule struct {
	fieldName M.StringId
	S.PropertyExpectation
}

//
type PendingRules []PropertyRule
