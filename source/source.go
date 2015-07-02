package source

import (
	"fmt"
	E "github.com/ionous/sashimi/event"
	G "github.com/ionous/sashimi/game" // for callbacks
	"github.com/ionous/sashimi/util/errutil"
)

//
// Placeholder for information about the location of definitions
// the runtime might implement an interface for this, or we could use stringer
// (ex. for handling compile time or run time errors )
//
type Code string

func (this Code) Errorf(format string, a ...interface{}) error {
	return errutil.Func(func() string {
		s := fmt.Errorf(format, a...)
		return fmt.Sprintf("Error (%s): %s", this, s)
	})
}

//
type IStatement interface {
	Source() Code // placeholder
}

//
//
//
type Blocks struct {
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
	blocks Blocks
}

func (this *BuildingBlocks) GetBlocks() Blocks {
	return this.blocks
}

//
func (this *BuildingBlocks) NewActionAssertion(
	actionName string,
	eventName string,
	source string,
	target string,
	context string,
) (err error) {
	fields := ActionAssertionFields{actionName, eventName, source, target, context}
	statement := ActionStatement{fields, ""}
	this.blocks.Actions = append(this.blocks.Actions, statement)
	return err
}

//
func (this *BuildingBlocks) NewActionHandler(
	owner string,
	action string,
	callback G.Callback,
	phase E.Phase,
) (err error) {
	fields := RunFields{owner, action, callback, phase}
	statement := RunStatement{fields, ""}
	this.blocks.ActionHandlers = append(this.blocks.ActionHandlers, statement)
	return err
}

func (this *BuildingBlocks) NewAlias(
	key string,
	phrases []string,
) (err error) {
	a := AliasStatement{key, phrases, ""}
	this.blocks.Aliases = append(this.blocks.Aliases, a)
	return err
}

//
func (this *BuildingBlocks) NewAssertion(
	owner string,
	called string,
	opts map[string]string,
) (err error) {
	a := AssertionStatement{owner, called, opts}
	this.blocks.Asserts = append(this.blocks.Asserts, a)
	return err
}

//
func (this *BuildingBlocks) NewChoice(fields ChoiceFields, source Code,
) (err error) {
	choice := ChoiceStatement{fields, source}
	this.blocks.Choices = append(this.blocks.Choices, choice)
	return err
}

//
func (this *BuildingBlocks) NewEnumeration(fields EnumFields, source Code,
) (err error) {
	enum := EnumStatement{fields, source}
	this.blocks.Enums = append(this.blocks.Enums, enum)
	return err
}

//
func (this *BuildingBlocks) NewEventHandler(fields ListenFields, source Code,
) (err error) {
	statement := ListenStatement{fields, source}
	this.blocks.EventHandlers = append(this.blocks.EventHandlers, statement)
	return err
}

//
func (this *BuildingBlocks) NewKeyValue(fields KeyValueFields, source Code,
) (err error) {
	kv := KeyValueStatement{fields, source}
	this.blocks.KeyValues = append(this.blocks.KeyValues, kv)
	return err
}

//
func (this *BuildingBlocks) NewMultiValue(fields MultiValueFields, source Code,
) (err error) {
	mv := MultiValueStatement{fields, source}
	this.blocks.MultiValues = append(this.blocks.MultiValues, mv)
	return err
}

//
func (this *BuildingBlocks) NewProperty(fields PropertyFields, source Code,
) (err error) {
	prop := PropertyStatement{fields, source}
	this.blocks.Properties = append(this.blocks.Properties, prop)
	return err
}

//
func (this *BuildingBlocks) NewRelative(fields RelativeFields, source Code,
) (err error) {
	rel := RelativeStatement{fields, source}
	this.blocks.Relatives = append(this.blocks.Relatives, rel)
	return err
}
