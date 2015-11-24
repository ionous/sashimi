package model

import "github.com/ionous/sashimi/util/ident"

// An indexed set of values.
type EnumModel struct {
	Choices []ident.Id `json:"choices"`
}

// 1 based
func (e *EnumModel) ChoiceToIndex(choice ident.Id) (ret int) {
	for i, c := range e.Choices {
		if choice == c {
			ret = i + 1
			break
		}
	}
	return
}
func (e *EnumModel) IndexToChoice(idx int) (ret ident.Id) {
	if idx > 0 && idx <= len(e.Choices) {
		ret = e.Choices[idx-1]
	}
	return
}

// it's probably best for constraints to live elsewhere
// so that the choices arent being duplicated at every property
// we would need to distill the contraints though for each new lvl.
//
// EnumConstraint// type EnumConstraint struct {
// 	Only         ident.Id
// 	Never        []ident.Id
// 	Usual        ident.Id
// 	UsuallyLocal bool // usual set for cons constraint, or for some ancestor?
// }

func (e *EnumModel) Best() (ret ident.Id) {
	return e.Choices[0]
	// 	switch {
	// 	case !e.Only.Empty():
	// 		ret = e.Only

	// 	case !e.Usual.Empty():
	// 		ret = e.Usual

	// 	default:
	// 		for _, choice := range e.Choices {
	// 			if !e.IsNever(choice) {
	// 				ret = choice
	// 				break
	// 			}
	// 		}
	// 	}
	// 	return
}

// func (e *EnumModel) IsNever(id ident.Id) (ret bool) {
// 	for _, n := range e.Never {
// 		if n == id {
// 			ret = true
// 			break
// 		}
// 	}
// 	return
// }
