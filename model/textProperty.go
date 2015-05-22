package model

type TextProperty struct {
	id   StringId
	name string
}

// FIX?  it might be nicer if the model had a builder with an interface for all such new(s)
func NewTextProperty(id StringId, name string) *TextProperty {
	return &TextProperty{id, name}
}

func (this *TextProperty) Id() StringId {
	return this.id
}

func (this *TextProperty) Name() string {
	return this.name
}
