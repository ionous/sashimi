package source

// Code records the origin of statements.
type Code string

func (c Code) String() string {
	return string(c)
}

// IStatement provides a uniform way of locating user script code.
type IStatement interface {
	Source() Code
}

// Statements contains all possible user script code.
type Statements struct {
	ActionHandlers []RunStatement
	Aliases        []AliasStatement
	Asserts        []AssertionStatement
	Actions        []ActionStatement
	Choices        []ChoiceStatement
	Enums          []EnumStatement
	EventHandlers  []ListenStatement
	KeyValues      []KeyValueStatement
	MultiValues    []MultiValueStatement
	Properties     []PropertyStatement
	Relatives      []RelativeStatement
}

// UnknownLocation is a stand-in for the file and line of the phrase used to build a statement. MARS: remove this and replace with the proper file and line!
const UnknownLocation = Code("unknown")

//
func (s *Statements) NewActionAssertion(
	fields ActionAssertionFields,
	source Code,
) (err error) {
	statement := ActionStatement{fields, source}
	s.Actions = append(s.Actions, statement)
	return err
}

//
func (s *Statements) NewActionHandler(fields RunFields, source Code,
) (err error) {
	statement := RunStatement{fields, source}
	s.ActionHandlers = append(s.ActionHandlers, statement)
	return err
}

//
func (s *Statements) NewAlias(fields AliasFields, source Code,
) (err error) {
	a := AliasStatement{fields, source}
	s.Aliases = append(s.Aliases, a)
	return err
}

//
func (s *Statements) NewAssertion(fields AssertionFields, source Code,
) (err error) {
	a := AssertionStatement{fields, source}
	s.Asserts = append(s.Asserts, a)
	return err
}

//
func (s *Statements) NewChoice(fields ChoiceFields, source Code,
) (err error) {
	choice := ChoiceStatement{fields, source}
	s.Choices = append(s.Choices, choice)
	return err
}

//
func (s *Statements) NewEnumeration(fields EnumFields, source Code,
) (err error) {
	enum := EnumStatement{fields, source}
	s.Enums = append(s.Enums, enum)
	return err
}

//
func (s *Statements) NewEventHandler(fields ListenFields, source Code,
) (err error) {
	statement := ListenStatement{fields, source}
	s.EventHandlers = append(s.EventHandlers, statement)
	return err
}

//
func (s *Statements) NewKeyValue(fields KeyValueFields, source Code,
) (err error) {
	kv := KeyValueStatement{fields, source}
	s.KeyValues = append(s.KeyValues, kv)
	return err
}

//
func (s *Statements) NewMultiValue(fields MultiValueFields, source Code,
) (err error) {
	mv := MultiValueStatement{fields, source}
	s.MultiValues = append(s.MultiValues, mv)
	return err
}

//
func (s *Statements) NewProperty(fields PropertyFields, source Code,
) (err error) {
	prop := PropertyStatement{fields, source}
	s.Properties = append(s.Properties, prop)
	return err
}

//
func (s *Statements) NewRelative(fields RelativeProperty, source Code,
) (err error) {
	rel := RelativeStatement{fields, source}
	s.Relatives = append(s.Relatives, rel)
	return err
}
