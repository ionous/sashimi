package metal

import "github.com/ionous/sashimi/util/ident"

// uses pointers re: gopherjs
type enumValue struct{ PanicValue }

func (p *enumValue) GetState() (ret ident.Id) {
	return p.getId()
}

// FIX: constraints
func (p *enumValue) SetState(c ident.Id) (err error) {
	return p.SetGeneric(c)
}
