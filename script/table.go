package script

import (
	S "github.com/ionous/sashimi/source"
)

// Table statements create, or add data to, one or more instances at a time.
// The columns must match existing properties, noting that all objects have an implicit "name" property.
// If a table lacks a "name" column, the compiler will generate a unqiue name automatically.
func Table(columns ...string) TableFragment {
	return TableFragment{NewOrigin(2), columns, nil}
}

// TableFragment assists in the generation of
type TableFragment struct {
	origin  Origin
	columns []string
	rows    []S.Row
}

// note: was "contains", changed to "Has" to align text with "And" statements
func (this TableFragment) Has(values ...interface{}) TableAndFragment {
	this.rows = append(this.rows, values)
	return TableAndFragment{this}
}

//
type TableAndFragment struct {
	TableFragment
}

//
func (this TableAndFragment) And(values ...interface{}) TableAndFragment {
	this.rows = append(this.rows, values)
	return this
}

//
func (this TableFragment) MakeStatement(b SubjectBlock) (err error) {
	fields := S.MultiValueFields{b.subject, this.columns, this.rows}
	return b.NewMultiValue(fields, this.origin.Code())
}
