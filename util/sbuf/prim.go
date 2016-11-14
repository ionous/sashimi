package sbuf

import (
	"reflect"
	"strconv"
)

type Bool struct {
	Bool bool
}
type Type struct {
	Value interface{}
}
type Value struct {
	Value interface{}
}
type String struct {
	S string
}
type Str struct {
	S interface{}
}
type Rune struct {
	Rune rune
}
type Int struct {
	Int int
}
type Int64 struct {
	Int int64
}
type Float struct {
	Float float64
}
type Error struct {
	Error error
}
type Hex64 struct {
	Uint uint64
}
type Quote struct {
	Value      Switch
	Start, End string
}

func Q(args ...interface{}) Quote {
	return Quote{args, "'", "'"}
}
func P(args ...interface{}) Quote {
	return Quote{args, "(", ")"}
}

func (op Bool) String() string {
	return strconv.FormatBool(op.Bool)
}
func (op String) String() string {
	return op.S
}
func (op Str) String() (ret string) {
	if s, ok := op.S.(Stringer); !ok {
		ret = "###"
	} else {
		ret = s.String()
	}
	return
}
func (op Rune) String() string {
	return string(op.Rune)
}
func (op Int) String() string {
	return formatInt(int64(op.Int))
}
func (op Int64) String() string {
	return formatInt(op.Int)
}
func (op Float) String() string {
	return formatFloat(op.Float)
}
func (op Error) String() string {
	return op.Error.Error()
}
func (op Hex64) String() string {
	return formatUint(op.Uint)
}
func (op Type) String() (ret string) {
	if op.Value == nil {
		ret = "<nil>"
	} else {
		ret = reflect.TypeOf(op.Value).String()
	}
	return
}
func (op Value) String() (ret string) {
	if op.Value == nil {
		ret = "<nil>"
	} else {
		ret = reflect.ValueOf(op.Value).String()
	}
	return
}
func (op Quote) String() (ret string) {
	return op.Start + op.Value.String() + op.End
}

func formatFloat(v float64) string {
	return strconv.FormatFloat(float64(v), 'g', -1, 64)
}
func formatInt(v int64) string {
	return strconv.FormatInt(v, 10)
}
func formatUint(v uint64) string {
	return strconv.FormatUint(v, 16)
}
