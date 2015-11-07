package api

import "github.com/ionous/sashimi/util/ident"

type Model interface {
	NumAction() int
	ActionNum(int) Action
	GetAction(ident.Id) (Action, bool)

	NumEvent() int
	EventNum(int) Event
	GetEvent(ident.Id) (Event, bool)

	NumClass() int
	ClassNum(int) Class
	GetClass(ident.Id) (Class, bool)

	NumInstance() int
	InstanceNum(int) Instance
	GetInstance(ident.Id) (Instance, bool)

	NumParserAction() int
	ParserActionNum(int) ParserAction

	AreCompatible(child, parent ident.Id) bool
	Pluralize(string) string
}

type Action interface {
	GetId() ident.Id
	// GetActionName returns the original name given by the scripter.
	GetActionName() string
	// GetEvent: raised by this action when the action occurs.
	GetEvent() Event
	// GetNouns: the classes for required by the action.
	GetNouns() Nouns
}

type Event interface {
	GetId() ident.Id
	// GetEventName returns the original name given by the scripter.
	GetEventName() string
}

type Prototype interface {
	GetId() ident.Id
	// GetParentClass returns nil for classes if no parent;
	// panics if no class can be found for an instnace.
	GetParentClass() Class
	GetOriginalName() string

	NumProperty() int
	PropertyNum(int) Property
	GetProperty(ident.Id) (Property, bool)

	// GetPropertyByChoice evalutes all properties to find an enumeration which can store the passed choice
	GetPropertyByChoice(choice ident.Id) (Property, bool)
}

type Instance Prototype
type Class Prototype

type PropertyType int

const (
	InvalidProperty PropertyType = iota
	NumProperty                  // float32
	TextProperty                 // string
	StateProperty                // string.Id
	ObjectProperty
	ArrayProperty = 1 << 16
)

type Property interface {
	GetId() ident.Id
	GetType() PropertyType
	//GetObjectType()?
	// or maybe IsCompatible(inst) bool
	GetValue() Value
	GetValues() Values
}

type Values interface {
	NumValue() int
	ValueNum(int) Value

	ClearValues()

	AppendNum(float32)
	AppendText(string)
	AppendObject(ident.Id) error

	// RemoveValue(int)
}

type Value interface {
	GetNum() float32
	SetNum(float32)

	GetText() string
	SetText(string)

	GetState() ident.Id
	SetState(ident.Id)

	// FIX : Relations relate Objects -> instances
	GetObject() ident.Id
	SetObject(ident.Id) error
}

// NOTE: ParserActions aren't id'd so, they are represented as structs.
type ParserAction struct {
	Action
	Commands []string
}

type NounType int

const (
	SourceNoun NounType = iota
	TargetNoun
	ContextNoun
)

// a list of nouns
type Nouns []ident.Id

func (n Nouns) GetNounCount() int {
	return len(n)
}

func (n Nouns) Get(t NounType) (ret ident.Id) {
	if int(t) < len(n) {
		ret = n[t]
	}
	return ret
}
