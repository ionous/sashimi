package compiler

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
)

//
// A property-pair declaration which generates a model.Relation.
//
type PendingRelation struct {
	name     string           // name of the relation
	src, dst *PendingRelative // copied from the class after verifying the two sides are compatible
}

//
// Assign a relative to the relation.
// There can be at most two relatives per relation, and both sides must relatesTo each other without conflict.
//
func (this *PendingRelation) setRelative(pending PendingRelative) (err error) {
	// setup name:
	name := pending.viaRelation
	if this.name == "" {
		this.name = name
	} else if this.name != name {
		err = fmt.Errorf("relation names don't match %s(old) %s(new)", this.name, name)
	}
	// setup relations:
	if err == nil {
		if !pending.isRev {
			if this.src != nil {
				err = fmt.Errorf("src already set %s %s", this.src, pending)
			} else if this.dst != nil && this.dst.relatesTo != pending.class {
				err = fmt.Errorf("src doesn't match dst '%s' '%s' in '%s'", this.src, pending, name)
			} else {
				this.src = &pending
			}
		} else {
			if this.dst != nil {
				err = fmt.Errorf("dst already set %s %s", this.dst, pending)
			} else if this.src != nil && this.src.relatesTo != pending.class {
				err = fmt.Errorf("dst doesn't match src '%s' '%s' in '%s'", this.dst, pending, name)
			} else {
				this.dst = &pending
			}
		}
	}
	return err
}

//
// Generate the model Relation from this PendingRelation.
//
func (this PendingRelation) makeRelation(id M.StringId) (rel M.Relation, err error) {
	if this.src == nil || this.dst == nil {
		err = fmt.Errorf("missing half of relation %v", this)
	} else {
		src := M.HalfRelation{this.src.class.id, this.src.id}
		dst := M.HalfRelation{this.dst.class.id, this.dst.id}
		rel = M.NewRelation(id, this.name, src, dst, this.style())
	}
	return rel, err
}

//
// deduce the relationship style based on how how many dst are pointed to by src, and vice versa
//
func (this PendingRelation) style() (style M.RelationStyle) {
	if this.src.toMany {
		if this.dst.toMany {
			style = M.ManyToMany
		} else {
			style = M.OneToMany // been drinking again, eh?
		}
	} else {
		if this.dst.toMany {
			style = M.ManyToOne
		} else {
			style = M.OneToOne
		}
	}
	return style
}
