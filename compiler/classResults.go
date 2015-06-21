package compiler

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
	S "github.com/ionous/sashimi/source"
)

type ClassResult struct {
	class *M.ClassInfo
	err   error
}

type ClassResultMap map[*PendingClass]*ClassResult

type ClassResults struct {
	pending   PendingClasses
	results   ClassResultMap
	relatives *RelativeFactory
}

//
//
//
func newResults(classes PendingClasses, relatives *RelativeFactory) ClassResults {
	results := make(ClassResultMap)
	// add a placeholder for the null/root class
	results[nil] = &ClassResult{nil, nil}
	return ClassResults{classes, results, relatives}
}

//
// exactly duplicated in makeInstances()
//
func (this ClassResults) finalizeClasses() (
	classes M.ClassMap,
	err error,
) {
	classes = make(M.ClassMap)

	for name, pending := range this.pending {
		info, e := this.makeClass(pending)
		if e != nil {
			err = AppendError(err, e)
		} else {
			classes[name] = info
		}
	}
	return classes, err
}

//
// make a class and its parents
//
func (this ClassResults) makeClass(pending *PendingClass,
) (cls *M.ClassInfo, err error) {
	cr := this.results[pending]
	if cr != nil {
		cls, err = cr.class, cr.err
	} else {
		cls, err = this._makeClass(pending)
		result := &ClassResult{cls, err}
		this.results[pending] = result
	}
	return cls, err
}

//
// recusively make class and parents
func (this ClassResults) _makeClass(pending *PendingClass,
) (cls *M.ClassInfo, err error) {
	// both so that we can have a direct pointer to our parent,
	// and so that we can find parent properties for adding constraints
	parent, e := this.makeClass(pending.parent)
	if e != nil {
		err = e
	} else {
		// make all the simple properties
		props := pending.makePropertySet()

		// distill all rules
		constraints := make(M.ConstraintSet)

		for _, rule := range pending.rules {
			// find prop for the rule
			prop, propFound := props[rule.fieldName]
			if !propFound && parent != nil {
				prop, propFound = parent.PropertyById(rule.fieldName)
			}

			// apply rule to property
			switch prop := prop.(type) {
			case nil:
				e := fmt.Errorf("rule specified for unknown field %s", rule.fieldName)
				err = AppendError(err, e)

			default:
				e := fmt.Errorf("rule specified for non-enum field %s", rule.fieldName)
				err = AppendError(err, e)

			case *M.EnumProperty:
				// find a constraint for the rule
				cons := constraints[rule.fieldName]
				if cons == nil && parent != nil {
					pcons := parent.ConstraintById(rule.fieldName)
					if pcons != nil {
						cons = pcons.Copy()
						constraints[rule.fieldName] = cons
					}
				}
				if cons == nil {
					cons = M.NewConstraint(prop.Enumeration)
					constraints[rule.fieldName] = cons
				}
				// add the new rule
				switch rule.RuleType() {
				case S.AlwaysExpect:
					e := cons.Always(rule.RuleValue())
					err = AppendError(err, e)
				case S.UsuallyExpect:
					e := cons.Usually(rule.RuleValue())
					err = AppendError(err, e)
				case S.SeldomExpect:
					e := cons.Seldom(rule.RuleValue())
					err = AppendError(err, e)
				case S.NeverExpect:
					e := cons.Exclude(rule.RuleValue())
					err = AppendError(err, e)
				default:
					e := fmt.Errorf("unknown type for rule %v", rule)
					err = AppendError(err, e)
				} // end rule switch
			} // end prop switch
		} // end rule loop

		// add relation properties
		for id, r := range pending.relatives {
			// get-or-create the associated relationship:
			if prop, e := this.relatives.makeRelative(r.rel); e != nil {
				err = AppendError(err, e)
			} else {
				props[id] = prop
			}
		}
		if err == nil {
			cls = M.NewClassInfo(
				parent,
				pending.id,
				pending.name,
				pending.singular,
				props,
				constraints)
		}
	}
	return cls, err
}
