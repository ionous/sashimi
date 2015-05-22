package model

type RelationStyle string

const (
	OneToOne   RelationStyle = "OneToOne"
	OneToMany                = "OneToMany"
	ManyToOne                = "ManyToOne"
	ManyToMany               = "ManyToMany"
)

//
// re: string ids.
//
// it seems be nice to have direct pointers in the model layer;
// but, it's also nice to have classes "sealed" once they are created;
// yet, relation properties can cause cycles, in some cases even pointing back to their own class.
//
// it's really not clear what's best, other than eventually, it needs to be consistent.
// some options are:
//
// 1) move creation into this package via a model maker API;
// classes, etc. are only partial during build --
// the final Info api is only exposed externally once its all finished.
//
// 2) instead of pointers, the models could reference only by id -- typed or otherwise;
// this would make them more closely match data on disk, db
//
// 3) some "core" class object -- +/- non-relation properties -- could be layered underneath C+R properties
//
type Relation struct {
	id       StringId
	name     string
	src, dst StringId // classes
	style    RelationStyle
}

func (this Relation) Id() StringId {
	return this.id
}

func (this Relation) Name() string {
	return this.name
}

func (this Relation) Source() StringId {
	return this.src
}

func (this Relation) Destination() StringId {
	return this.dst
}

func (this Relation) Style() RelationStyle {
	return this.style
}

// func (this Relation) Other(test StringId) (ret StringId, many bool) {
// 	if this.src == test {
// 		ret = this.dst
// 	} else {
// 		ret = this.src
// 	}
// 	return
// }

type RelationMap map[StringId]Relation

func NewRelation(id StringId, name string, src, dst StringId, style RelationStyle) Relation {
	return Relation{id, name, src, dst, style}
}
