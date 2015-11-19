package model

import (
	"github.com/ionous/sashimi/util/ident"
	"strings"
)

type ClassModel struct {
	Id         ident.Id        `json:"id"`
	Parents    []ident.Id      `json:"parents"`
	Plural     string          `json:"plural"`
	Singular   string          `json:"singular"`
	Properties []PropertyModel `json:"properties"`
}

//
func (cls ClassModel) String() string {
	return cls.Plural
}

func (cls ClassModel) Parent() (ret ident.Id) {
	if len(cls.Parents) > 0 {
		ret = cls.Parents[0]
	}
	return
}

// FindProperty does not search parent classes.
func (cls ClassModel) FindPropertyById(id ident.Id) (ret *PropertyModel, okay bool) {
	for i, p := range cls.Properties {
		if p.Id == id {
			ret, okay = &cls.Properties[i], true
			break
		}
	}
	return
}

func (cls ClassModel) FindProperty(name string) (ret *PropertyModel, okay bool) {
	for i, p := range cls.Properties {
		if strings.EqualFold(p.Name, name) {
			ret, okay = &cls.Properties[i], true
			break
		}
	}
	return
}
