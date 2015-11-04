package internal

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/util/ident"
)

// PendingRelation: a property-pair declaration which generates a model.Relation.
type PendingRelation struct {
	name     string
	src, dst PendingRelative
}

// PendingRelative: copied from the class after verifying the two sides are compatible
type PendingRelative struct {
	Class ident.Id
	M.RelativeProperty
}

func (r PendingRelative) Empty() bool {
	return r.Class.Empty()
}

// PendingRelations maps relation-id to a pending relation.
type PendingRelations map[ident.Id]PendingRelation

// setRelative assigns the side ( src/ dst ) of a relation.
// There can be at most two relatives per relation, and both sides must relate to each other without conflict.
func (p *PendingRelation) setRelative(name string, pending PendingRelative) (err error) {
	// setup name:
	if p.name == "" {
		p.name = name
	} else if p.name != name {
		err = fmt.Errorf("relation names don't match %s(old) %s(new)", p.name, name)
	}
	// setup relations:
	if !pending.IsRev {
		err = setRelative(&p.src, name, p.dst, pending)
	} else {
		err = setRelative(&p.dst, name, p.src, pending)
	}
	return err
}

func setRelative(src *PendingRelative, name string, dst, pending PendingRelative) (err error) {
	if !src.Empty() {
		err = fmt.Errorf("src already set %s %+v %+v", name, src, pending)
	} else if !dst.Empty() && dst.Relates != pending.Class {
		err = fmt.Errorf("src doesn't match dst '%s' '%s' in '%s'", src, pending, name)
	} else {
		*src = pending
	}
	return err
}

// makeRelation generates the model Relation from p PendingRelation.
func (p PendingRelation) makeRelation(id ident.Id) (rel M.Relation, err error) {
	if p.src.Empty() || p.dst.Empty() {
		err = fmt.Errorf("missing half of relation %v", p)
	} else {
		src := M.HalfRelation{p.src.Class, p.src.Id}
		dst := M.HalfRelation{p.dst.Class, p.dst.Id}
		rel = M.NewRelation(id, p.name, src, dst, p.style())
	}
	return rel, err
}

// style deduces RelationStyle based on how how many dst are pointed to by src, and vice versa
func (p PendingRelation) style() (style M.RelationStyle) {
	if p.src.IsMany {
		if p.dst.IsMany {
			style = M.ManyToMany
		} else {
			style = M.OneToMany // been drinking again, eh?
		}
	} else {
		if p.dst.IsMany {
			style = M.ManyToOne
		} else {
			style = M.OneToOne
		}
	}
	return style
}
