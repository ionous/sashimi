package compiler

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
)

// map future relation string id to data about the relation
type RelativeFactory struct {
	names     NameScope
	relations PendingRelations
}
type PendingRelations map[M.StringId]PendingRelation

func newRelativeFactory(names NameScope) *RelativeFactory {
	return &RelativeFactory{names, make(PendingRelations)}
}

// make a relative property by, interally, generating a relation
func (this *RelativeFactory) makeRelative(propid M.StringId, pending *PendingClass, prop PendingRelative) (ret *M.RelativeProperty, err error) {
	// sort out the source and destination names
	srcName, dstName := pending.id.String(), prop.relatesTo.id.String()
	if prop.isRev {
		dstName, srcName = srcName, dstName
	}
	// create an id to match the two classes
	if relId, e := this.names.addName(prop.viaRelation, srcName+dstName); e != nil {
		err = e
	} else {
		// merge this relative with it's counter-part ( if any so far )
		rel := this.relations[relId]
		if e := rel.setName(prop.viaRelation); e != nil {
			err = e
		} else {
			// sort out the source and destination data in the relation
			src, dst := &rel.src, &rel.dst
			if prop.isRev {
				dst, src = src, dst
			}

			if e := src.setRelationTo(prop.relatesTo, prop.toMany); e != nil {
				err = e
			} else if e := dst.setRelationFrom(pending); e != nil {
				err = e
			} else {
				// write merged data back
				this.relations[relId] = rel

				// finally, create the relative property pointing to the generated relation data
				fields := M.RelativeFields{propid, prop.name, relId, src.class.id, prop.isRev, prop.toMany}
				ret = M.NewRelative(fields)
			}
		}
	}
	return ret, err
}

// every relation is composed to two sides, possibly gathered over time
type HalfRelation struct {
	class     *PendingClass // what do we point to
	toMany    bool
	toManySet bool
}

//
type PendingRelation struct {
	name string
	src  HalfRelation // what does the source container point to
	dst  HalfRelation // what does the dest container point to
	// and... relation by values? a class? a property?
}

//
func (this *PendingRelation) setName(name string) (err error) {
	if this.name == "" {
		this.name = name
	} else if this.name != name {
		err = fmt.Errorf("internal error? relation names don't match %s(old) %s(new)", this.name, name)
	}
	return err
}

//
func (this PendingRelation) makeRelation(id M.StringId) (rel M.Relation, err error) {
	if this.src.class == nil || this.dst.class == nil {
		err = fmt.Errorf("missing source half of relation", this)
	} else {
		src, dst := this.src.class.id, this.dst.class.id
		// so, this is a little confused:
		// up to this point, source is what the source points to
		// but the relation makes more sense when we talk about the source container
		// which is stored on the dest ( opposite ) side;
		// the style is correct, it still talks about how does source container point to the dest container
		rel = M.NewRelation(id, this.name, dst, src, this.style())
	}
	return rel, err
}

//
// deduce the relationship style based on how how many dst are pointed to by src, and vice versa
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

//
// from a single property we can't determine the full state of the relation's (reverse) mapping
// so, we only set the class; and leave toMany for the other side ( if it exists )
// the default reverse toMany is false.
func (this *HalfRelation) setRelationFrom(class *PendingClass) (err error) {
	return this._setClass(class)
}

//
func (this *HalfRelation) setRelationTo(class *PendingClass, toMany bool) (err error) {
	if e := this._setClass(class); e != nil {
		err = e
	} else {
		if !this.toManySet {
			this.toMany, this.toManySet = toMany, true
		} else if this.toMany != toMany {
			err = fmt.Errorf("relation doesn't match %s(was), %s;%b(to)", this, class, toMany)
		}
	}
	return err
}

//
func (this *HalfRelation) _setClass(class *PendingClass) (err error) {
	if this.class == nil {
		this.class = class
	} else if this.class != class {
		err = fmt.Errorf("relation doesn't match %s(was), %s(to)", this.class, class)
	}
	return err
}
