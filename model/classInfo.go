package model

type ClassInfo struct {
	parent      *ClassInfo
	id          StringId
	name        string
	singular    string
	props       PropertySet // properties only for this class
	constraints ConstraintSet
}

func NewClassInfo(
	parent *ClassInfo,
	id StringId,
	plural string,
	singular string,
	props PropertySet,
	constraints ConstraintSet,
) (cls *ClassInfo) {
	return &ClassInfo{parent, id, plural, singular, props, constraints}
}

//
func (this *ClassInfo) Id() StringId {
	return this.id
}

//
func (this *ClassInfo) Name() string {
	return this.name
}

//
func (this *ClassInfo) Singular() string {
	return this.singular
}

//
func (this *ClassInfo) String() string {
	return this.name
}

//
func (this *ClassInfo) Parent() *ClassInfo {
	return this.parent
}

//
func (this *ClassInfo) Properties() PropertySet {
	return this.props
}

//
// Returns a new property set consisting of all properties in this class and all parents
//
func (this *ClassInfo) AllProperties() PropertySet {
	props := make(PropertySet)
	this._flatten(props)
	return props
}

//
//
//
func (this *ClassInfo) FindProperty(name string) (prop IProperty, okay bool) {
	id := MakeStringId(name)
	return this.PropertyById(id)
}

//
//
//
func (this *ClassInfo) PropertyById(id StringId) (prop IProperty, okay bool) {
	prop, okay = this.props[id]
	if !okay && this.parent != nil {
		prop, okay = this.parent.PropertyById(id)
	}
	return prop, okay
}

//
//
//
func (this *ClassInfo) Constraints() ConstraintSet {
	return this.constraints
}

//
//
//
func (this *ClassInfo) Constraint(name string) IConstrain {
	id := MakeStringId(name)
	return this.ConstraintById(id)
}

//
//
//
func (this *ClassInfo) ConstraintById(id StringId) IConstrain {
	constraint := this.constraints[id]
	if constraint == nil && this.parent != nil {
		constraint = this.parent.ConstraintById(id)
	}
	return constraint
}

//
//
//
func (this *ClassInfo) HasParent(p *ClassInfo) (yes bool) {
	for c := this.Parent(); c != nil; c = c.Parent() {
		if c == p {
			yes = true
			break
		}
	}
	return yes
}

//
//
//
func (this *ClassInfo) PropertyByChoice(choice string) (
	prop *EnumProperty,
	index int,
	ok bool,
) {
	choiceId := MakeStringId(choice)
	prop, index = this._propertyByChoice(choiceId)
	return prop, index, prop != nil
}

func (this *ClassInfo) _propertyByChoice(choice StringId) (
	prop *EnumProperty,
	index int,
) {
	for _, p := range this.props {
		switch enum := p.(type) {
		case *EnumProperty:
			idx, err := enum.ChoiceToIndex(choice)
			if err == nil {
				prop = enum
				index = idx
			}
		}
		if prop != nil {
			break
		}
	}
	if prop == nil && this.parent != nil {
		prop, index = this.parent._propertyByChoice(choice)
	}
	return prop, index
}

// NOTE: does NOT check for conflicts.
// ( trying to be a little looser than normal,
// and get to the point where the model is known to be safe at creation time. )
func (this *ClassInfo) _flatten(props PropertySet) {
	if this.parent != nil {
		this.parent._flatten(props)
	}
	for k, prop := range this.props {
		props[k] = prop
	}
}
