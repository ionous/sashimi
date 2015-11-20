package metal

import (
	//"fmt"
	"github.com/ionous/sashimi/util/ident"
)

type enumValue struct{ panicValue }

func (p enumValue) GetState() (ret ident.Id) {
	// if enum, ok := p.mdl.Enumerations[p.prop.Id]; !ok {
	// 	panic(fmt.Sprintf("internal error, couldnt find enumeration '%s.%s'", p.src, p.prop.Id))
	// } else {
	// v := p.get()
	// if idx, ok := v.(int); !ok {
	// 	panic(fmt.Sprintf("internal error, couldnt convert state to int '%s.%s' %v(%T)", p.src, p.prop.Id, v, v))
	// } else {
	// 	ret = enum.IndexToChoice(idx)
	// }
	//}
	return p.get().(ident.Id)
}

// FIX: constraints
func (p enumValue) SetState(c ident.Id) (err error) {
	// if enum, ok := p.mdl.Enumerations[p.prop.Id]; !ok {
	// 	panic(fmt.Sprintf("internal error, couldnt find enumeration '%s.%s'", p.src, p.prop.Id))
	// } else if idx := enum.ChoiceToIndex(c); idx <= 0 {
	// 	err = fmt.Errorf("no such choice %s in '%s.%s'", c, p.src, p.prop.Id)
	// } else {
	// 	p.set(idx)
	// }
	return p.set(c)
}
