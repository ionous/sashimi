package resource

//
// NewObject creation.
//
func NewObject(id, class string) *Object {
	return &Object{id, class, make(Dict), make(map[string]Relationship), make(Dict)}
}

//
// SetAttr sets the attribute "key" to "value".
//
func (obj *Object) SetAttr(key string, value interface{}) *Object {
	obj.Attributes[key] = value
	return obj
}

//
// SetMeta sets the metada "key" to "value".
//
func (obj *Object) SetMeta(key string, value interface{}) *Object {
	obj.Meta[key] = value
	return obj
}

//
// SetRel adds object relations.
// FUTURE: use a builder for metadata support, etc.
//
func (obj *Object) SetRel(key string, data interface{}) *Object {
	obj.Relationships[key] = Relationship{Data: data}
	return obj
}
