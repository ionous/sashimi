package model

type EnumProperty struct {
	id   StringId
	name string
	Enumeration
}

func NewEnumProperty(id StringId, name string, src Enumeration) *EnumProperty {
	return &EnumProperty{id, name, src}
}

func (this *EnumProperty) Id() StringId {
	return this.id
}

func (this *EnumProperty) Name() string {
	return this.name
}
