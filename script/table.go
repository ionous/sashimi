package script

import (
	S "github.com/ionous/sashimi/source"
)

//
// Statement to add a row of columns to a table/view.
// The column names must match existing property names.
//
func Table(columns ...string) TableFragment {
	return TableFragment{NewOrigin(2), columns, nil}
}

//
//
//
type TableFragment struct {
	origin  Origin
	columns []string
	rows    []S.Row
}

func (this TableFragment) Contains(values ...interface{}) TableAndFragment {
	this.rows = append(this.rows, values)
	return TableAndFragment{this}
}

type TableAndFragment struct {
	TableFragment
}

func (this TableAndFragment) And(values ...interface{}) TableAndFragment {
	this.rows = append(this.rows, values)
	return this
}

//
func (this TableFragment) MakeStatement(b SubjectBlock) (err error) {
	fields := S.MultiValueFields{b.subject, this.columns, this.rows}
	return b.NewMultiValue(fields, this.origin.Code())
}
