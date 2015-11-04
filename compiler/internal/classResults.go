package internal

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
	S "github.com/ionous/sashimi/source"
	"github.com/ionous/sashimi/util/errutil"
)

// Each class has an associated result
type ClassResult struct {
	class *M.ClassInfo
	err   error
}

// Each pending class gets turned into a result
type ClassResultMap map[*PendingClass]*ClassResult

//
type ClassResults struct {
	pending   PendingClasses
	results   ClassResultMap
	relatives *RelativeFactory
}

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
		if info, e := this.makeClass(pending); e != nil {
			err = errutil.Append(err, e)
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
	} else if props, e := pending.makePropertySet(); e != nil {
		err = e
	} else {
		// distill all rules
		constraints := make(M.ConstraintMap)

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
				err = errutil.Append(err, e)

			default:
				e := fmt.Errorf("rule specified for non-enum field %s", rule.fieldName)
				err = errutil.Append(err, e)

			case M.EnumProperty:
				// find a constraint for the rule
				var cons M.IConstrain
				found := false
				if c, ok := constraints[rule.fieldName]; ok {
					cons = c
					found = true
				} else if parent != nil {
					if p, ok := parent.Constraints.ConstraintById(rule.fieldName); ok {
						cons = p.Copy()
						constraints[rule.fieldName] = cons
						found = true
					}
				}
				if !found {
					cons = M.NewConstraint(prop.Enumeration)
					constraints[rule.fieldName] = cons
				}
				// add the new rule
				switch rule.RuleType() {
				case S.AlwaysExpect:
					e := cons.Always(rule.RuleValue())
					err = errutil.Append(err, e)
				case S.UsuallyExpect:
					e := cons.Usually(rule.RuleValue())
					err = errutil.Append(err, e)
				case S.SeldomExpect:
					e := cons.Seldom(rule.RuleValue())
					err = errutil.Append(err, e)
				case S.NeverExpect:
					e := cons.Exclude(rule.RuleValue())
					err = errutil.Append(err, e)
				default:
					e := fmt.Errorf("unknown type for rule %v", rule)
					err = errutil.Append(err, e)
				} // end rule switch
			} // end prop switch
		} // end rule loop

		if err == nil {
			var parentConstraints *M.ConstraintSet
			if parent != nil {
				parentConstraints = &parent.Constraints
			}
			cls = &M.ClassInfo{
				Parent:      parent,
				Id:          pending.id,
				Plural:      pending.name,
				Singular:    pending.singular,
				Properties:  props,
				Constraints: M.ConstraintSet{parentConstraints, constraints},
			}
		}
	}
	return cls, err
}
