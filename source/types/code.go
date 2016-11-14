package types

// Code records the origin of statements.
type Code string

func (c Code) String() string {
	return string(c)
}
