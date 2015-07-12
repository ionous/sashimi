package source

import (
	"fmt"
	E "github.com/ionous/sashimi/event"
	G "github.com/ionous/sashimi/game" // for callbacks
	"github.com/ionous/sashimi/util/errutil"
)

//
// Placeholder for information about the location of definitions
// the runtime might implement an interface for blocks, or we could use stringer
// (ex. for handling compile time or run time errors )
//
type Code string

func (code Code) Errorf(format string, a ...interface{}) error {
	return errutil.Func(func() string {
		s := fmt.Errorf(format, a...)
		return fmt.Sprintf("Error (%s): %s", code, s)
	})
}

//
type IStatement interface {
	Source() Code // placeholder
}

//
//
//
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

type BuildingBlocks struct {
	statements Statements
}

func (blocks *BuildingBlocks) GetStatements() Statements {
	return blocks.statements
}

//
func (blocks *BuildingBlocks) NewActionAssertion(
	actionName string,
	eventName string,
	source string,
	target string,
	context string,
) (err error) {
	fields := ActionAssertionFields{actionName, eventName, source, target, context}
	statement := ActionStatement{fields, ""}
	blocks.statements.Actions = append(blocks.statements.Actions, statement)
	return err
}

//
func (blocks *BuildingBlocks) NewActionHandler(
	owner string,
	action string,
	callback G.Callback,
	phase E.Phase,
) (err error) {
	fields := RunFields{owner, action, callback, phase}
	statement := RunStatement{fields, ""}
	blocks.statements.ActionHandlers = append(blocks.statements.ActionHandlers, statement)
	return err
}

func (blocks *BuildingBlocks) NewAlias(
	key string,
	phrases []string,
) (err error) {
	a := AliasStatement{key, phrases, ""}
	blocks.statements.Aliases = append(blocks.statements.Aliases, a)
	return err
}

//
func (blocks *BuildingBlocks) NewAssertion(fields AssertionFields, source Code,
) (err error) {
	a := AssertionStatement{fields, source}
	blocks.statements.Asserts = append(blocks.statements.Asserts, a)
	return err
}

//
func (blocks *BuildingBlocks) NewChoice(fields ChoiceFields, source Code,
) (err error) {
	choice := ChoiceStatement{fields, source}
	blocks.statements.Choices = append(blocks.statements.Choices, choice)
	return err
}

//
func (blocks *BuildingBlocks) NewEnumeration(fields EnumFields, source Code,
) (err error) {
	enum := EnumStatement{fields, source}
	blocks.statements.Enums = append(blocks.statements.Enums, enum)
	return err
}

//
func (blocks *BuildingBlocks) NewEventHandler(fields ListenFields, source Code,
) (err error) {
	statement := ListenStatement{fields, source}
	blocks.statements.EventHandlers = append(blocks.statements.EventHandlers, statement)
	return err
}

//
func (blocks *BuildingBlocks) NewKeyValue(fields KeyValueFields, source Code,
) (err error) {
	kv := KeyValueStatement{fields, source}
	blocks.statements.KeyValues = append(blocks.statements.KeyValues, kv)
	return err
}

//
func (blocks *BuildingBlocks) NewMultiValue(fields MultiValueFields, source Code,
) (err error) {
	mv := MultiValueStatement{fields, source}
	blocks.statements.MultiValues = append(blocks.statements.MultiValues, mv)
	return err
}

//
func (blocks *BuildingBlocks) NewProperty(fields PropertyFields, source Code,
) (err error) {
	prop := PropertyStatement{fields, source}
	blocks.statements.Properties = append(blocks.statements.Properties, prop)
	return err
}

//
func (blocks *BuildingBlocks) NewRelative(fields RelativeFields, source Code,
) (err error) {
	rel := RelativeStatement{fields, source}
	blocks.statements.Relatives = append(blocks.statements.Relatives, rel)
	return err
}
