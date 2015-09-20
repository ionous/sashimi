package resource

// Document provides json-data roughly following the description at http://jsonapi.org/format/#document-structure.
// DocumentBuilder provides an API for generating these structures properly.
type Document struct {
	// data can either be an Object or an array of Objects
	// (the lack of omitempty vs omitnil is super-annoying.)
	Data     interface{} `json:"data"`
	Errors   []Error     `json:"error,omitempty"`
	Meta     Dict        `json:"meta,omitempty"`
	Included []*Object   `json:"included,omitempty"`
	Links    Links       `json:"links,omitempty"`
	//jsonapi
}

type Links map[string]Link
type Link string

// for deserializing
type ObjectDocument struct {
	Data     Object   `json:"data"`
	Meta     Dict     `json:"meta,omitempty"`
	Included []Object `json:"included,omitempty"`
}

// for deserializing
type MultiDocument struct {
	Data     []Object `json:"data"`
	Meta     Dict     `json:"meta,omitempty"`
	Included []Object `json:"included,omitempty"`
}

type Object struct {
	Id            string                  `json:"id,omitempty"`
	Class         string                  `json:"type,omitempty"`
	Attributes    Dict                    `json:"attributes,omitempty"`
	Relationships map[string]Relationship `json:"relationships,omitempty"`
	Meta          Dict                    `json:"meta,omitempty"`
}

type Relationship struct {
	// data can either be an Object or an array of Objects
	// (the lack of omitempty vs omitnil is super-annoying.)
	Data interface{} `json:"data,omitempty"`
	Meta Dict        `json:"meta,omitempty"`
	// Self     Link    `json:"self,omitempty"`
	// Related  Link    `json:"related,omitempty"`
}

type Error struct {
	//Id string `json:"id,omitempty"`
	// status, code, title, detail,
	// links, source,
	Code string `json:"meta,omitempty"`
}
