package types

type NamedAction string
type NamedChoice string
type NamedChoices []string
type NamedClass string // or primitive
type NamedEvent string
type NamedEvents []string
type NamedNoun string // instance
type NamedProperty string
type NamedScript string
type NamedSubject string // class or instance
type PlayerInput []string

func (s NamedAction) String() string {
	return string(s)
}
func (s NamedChoice) String() string {
	return string(s)
}
func (s NamedChoices) Strings() []string {
	return []string(s)
}
func (s NamedClass) String() string {
	return string(s)
}
func (s NamedEvent) String() string {
	return string(s)
}
func (s NamedEvents) Strings() []string {
	return []string(s)
}
func (s NamedNoun) String() string {
	return string(s)
}
func (s NamedProperty) String() string {
	return string(s)
}
func (s NamedScript) String() string {
	return string(s)
}
func (s NamedSubject) String() string {
	return string(s)
}
func (s PlayerInput) Strings() []string {
	return []string(s)
}
