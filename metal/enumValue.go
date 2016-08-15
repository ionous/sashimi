package metal

import "github.com/ionous/sashimi/util/ident"

// uses pointers re: gopherjs
type enumValue struct{ panicValue }

func (p *enumValue) GetState() (ret ident.Id) {
	return p.getId()
}

// FIX: constraints
func (p *enumValue) SetState(c ident.Id) (err error) {
	return p.set(c)
}
