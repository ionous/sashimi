package resource

//
// object creation
//
func NewObject(id, class string) *Object {
	return &Object{id, class, make(Dict), make(map[string]Relationship), make(Dict)}
}

//
// Add an attribute to this object.
//
func (this *Object) SetAttr(key string, value interface{}) *Object {
	this.Attributes[key] = value
	return this
}

//
// Add object metadata.
//
func (this *Object) SetMeta(key string, value interface{}) *Object {
	this.Meta[key] = value
	return this
}

//
// Add object relations
// FUTURE: use a builder for metadata support, etc.
//
func (this *Object) SetRel(key string, data interface{}) *Object {
	this.Relationships[key] = Relationship{Data: data}
	return this
}
