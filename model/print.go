package model

import (
//	"fmt"
)

type printer func(...interface{})

func (model *Model) PrintModel(printer printer) {
	// if model == nil {
	// 	panic("missing results")
	// }
	// printer("*** Classes:")
	// for _, class := range model.Classes {
	// 	printer("\t", class.Plural, class.Singular)
	// 	cons := class.Constraints
	// 	printer("\t\t Constraints:")

	// 	if len(cons.Map) == 0 {
	// 		printer("\t\t\t (unconstrained)")
	// 	} else {
	// 		for field, con := range cons.Map {
	// 			printer("\t\t\t", fmt.Sprintf("%v: %T", field, con))
	// 		}
	// 	}

	// 	properties := class.Properties
	// 	printer("\t\t Fields:")
	// 	if len(properties) == 0 {
	// 		printer("\t\t\t (empty)")
	// 	} else {
	// 		for field, prop := range properties {
	// 			printer("\t\t\t", fmt.Sprintf("%v: %T", field, prop))
	// 			switch inner := prop.(type) {
	// 			case EnumProperty:
	// 				printer("\t\t\t ", inner.Values)

	// 			case RelativeProperty:
	// 				rel, _ := model.Relations[inner.Relation]
	// 				many := map[bool]string{true: "Many", false: "One"}[inner.IsMany]
	// 				printer("\t\t\t ", fmt.Sprintf("%s => '%s' ( '%s' )", many, inner.Relates, rel.Name))
	// 			}
	// 		}
	// 	}
	// }
	// printer("*** Relations:")
	// for id, rel := range model.Relations {
	// 	printer("\t", rel)
	// 	if table, ok := model.Tables[id]; ok {
	// 		printer("\t", table)
	// 	}
	// }
	// printer("*** Instances:")
	// for _, inst := range model.Instances {
	// 	printer("\t", inst)
	// 	all := inst.Class.AllProperties()
	// 	for id, _ := range all {
	// 		v, hadValue := inst.Values[id]
	// 		l := map[bool]string{false: "_", true: ""}[hadValue]
	// 		printer("\t\t", fmt.Sprintf("%v%v: %s", id, l, v))
	// 	}
	// }
	// printer("*** Actions:")
	// for _, act := range model.Actions {
	// 	source, target, context := act.Source(), act.Target(), act.Context()
	// 	printer("\t", act.ActionName, act.EventName, source, target, context)
	// }
}
