package compiler

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/util/ident"
)

//
// A property-pair declaration which generates a model.Relation.
//
type PendingRelation struct {
	name     string
	src, dst *M.RelativeProperty // copied from the class after verifying the two sides are compatible
}

//
// PendingRelations maps relation-id to a pending relation.
//
type PendingRelations map[ident.Id]PendingRelation

//
// Assign a relative to the relation.
// There can be at most two relatives per relation, and both sides must relate to each other without conflict.
//
func (this *PendingRelation) setRelative(name string, pending M.RelativeProperty) (err error) {
	// setup name:
	if this.name == "" {
		this.name = name
	} else if this.name != name {
		err = fmt.Errorf("relation names don't match %s(old) %s(new)", this.name, name)
	}

	// setup relations:
	if !pending.IsRev {
		if this.src != nil {
			err = fmt.Errorf("src already set %s %+v %+v", name, this.src, pending)
		} else if this.dst != nil && this.dst.Relates != pending.Class {
			err = fmt.Errorf("src doesn't match dst '%s' '%s' in '%s'", this.src, pending, name)
		} else {
			this.src = &pending
		}
	} else {
		if this.dst != nil {
			err = fmt.Errorf("dst already set %s %s", this.dst, pending)
		} else if this.src != nil && this.src.Relates != pending.Class {
			err = fmt.Errorf("dst doesn't match src '%s' '%s' in '%s'", this.dst, pending, name)
		} else {
			this.dst = &pending
		}
	}
	return err
}

//
// Generate the model Relation from this PendingRelation.
//
func (this PendingRelation) makeRelation(id ident.Id) (rel M.Relation, err error) {
	if this.src == nil || this.dst == nil {
		err = fmt.Errorf("missing half of relation %v", this)
	} else {
		src := M.HalfRelation{this.src.Class, this.src.Id}
		dst := M.HalfRelation{this.dst.Class, this.dst.Id}
		rel = M.NewRelation(id, this.name, src, dst, this.style())
	}
	return rel, err
}

//
// deduce the relationship style based on how how many dst are pointed to by src, and vice versa
//
func (this PendingRelation) style() (style M.RelationStyle) {
	if this.src.IsMany {
		if this.dst.IsMany {
			style = M.ManyToMany
		} else {
			style = M.OneToMany // been drinking again, eh?
		}
	} else {
		if this.dst.IsMany {
			style = M.ManyToOne
		} else {
			style = M.OneToOne
		}
	}
	return style
}
