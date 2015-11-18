package internal

import (
	"fmt"
	M "github.com/ionous/sashimi/compiler/xmodel"
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
	pending        PendingClasses
	results        ClassResultMap
	relatives      *RelativeFactory
	singleToPlural map[string]string
}

//
func newResults(classes PendingClasses, relatives *RelativeFactory) ClassResults {
	results := make(ClassResultMap)
	// add a placeholder for the null/root class
	results[nil] = &ClassResult{nil, nil}
	return ClassResults{classes, results, relatives, make(map[string]string)}
}

//
// exactly duplicated in makeInstances()
//
func (cr ClassResults) finalizeClasses() (
	classes M.ClassMap,
	singleToPlural map[string]string,
	err error,
) {
	classes = make(M.ClassMap)

	for name, pending := range cr.pending {
		if info, e := cr.makeClass(pending); e != nil {
			err = errutil.Append(err, e)
		} else {
			classes[name] = info
			singleToPlural = cr.singleToPlural
		}
	}
	return classes, singleToPlural, err
}

//
// make a class and its parents
//
func (cr ClassResults) makeClass(pending *PendingClass,
) (cls *M.ClassInfo, err error) {
	if res, ok := cr.results[pending]; ok {
		cls, err = res.class, res.err
	} else {
		cls, err = cr._makeClass(pending)
		result := &ClassResult{cls, err}
		cr.results[pending] = result
	}
	return cls, err
}

//
// recusively make class and parents
func (cr ClassResults) _makeClass(pending *PendingClass,
) (cls *M.ClassInfo, err error) {
	// both so that we can have a direct pointer to our parent,
	// and so that we can find parent properties for adding constraints
	parent, e := cr.makeClass(pending.parent)
	if e != nil {
		err = e
	} else if props, e := pending.makePropertySet(); e != nil {
		err = e
	} else {
		// distill all rules
		constraints := make(M.ConstraintMap)

		for _, rule := range pending.rules {
			// find prop for the rule
			prop, propFound := props[rule.propId]
			if !propFound && parent != nil {
				prop, propFound = parent.PropertyById(rule.propId)
			}

			// apply rule to property
			switch prop := prop.(type) {
			case nil:
				e := fmt.Errorf("rule specified for unknown field %s", rule.propId)
				err = errutil.Append(err, e)

			default:
				e := fmt.Errorf("rule specified for non-enum field %s", rule.propId)
				err = errutil.Append(err, e)

			case M.EnumProperty:
				// find a constraint for the rule
				var cons M.IConstrain
				found := false
				if c, ok := constraints[prop.GetId()]; ok {
					cons = c
					found = true
				} else if parent != nil {
					if p, ok := parent.Constraints.ConstraintById(prop.GetId()); ok {
						cons = p.Copy()
						constraints[prop.GetId()] = cons
						found = true
					}
				}
				if !found {
					cons = M.NewConstraint(prop.Enumeration)
					constraints[prop.GetId()] = cons
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
			cr.singleToPlural[pending.singular] = pending.name
		}
	}
	return cls, err
}
