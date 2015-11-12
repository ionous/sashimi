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

	NumRelation() int
	RelationNum(int) Relation
	GetRelation(ident.Id) (Relation, bool)

	NumParserAction() int
	ParserActionNum(int) ParserAction

	Pluralize(string) string
	AreCompatible(child, parent ident.Id) bool
	MatchNounName(string, func(ident.Id) bool) (tries int, err bool)
}

type Action interface {
	GetId() ident.Id
	// GetActionName returns the original name given by the scripter.
	GetActionName() string
	// GetEvent: raised by this action when the action occurs.
	GetEvent() Event
	// GetNouns: the classes for required by the action.
	GetNouns() Nouns
	//
	GetCallbacks() (Callbacks, bool)
}

type Callbacks interface {
	NumCallback() int
	CallbackNum(int) ident.Id
}

type Event interface {
	GetId() ident.Id
	// GetEventName returns the original name given by the scripter.
	GetEventName() string
	// GetAction: there shouldnt have to be a one to one mapping between actions ad events
	// that's just how it is right now :(
	GetAction() Action
	// GetListeners returns the capture or bubbling callbacks associated with this event
	// if GetListeners returns false, Listeners should be set to NoListeners.
	GetListeners(capture bool) (Listeners, bool)
}

type Listeners interface {
	NumListener() int
	ListenerNum(int) Listener
}

type NoListeners [0]Listener

func (no NoListeners) NumListener() int {
	return len(no)
}

func (no NoListeners) ListenerNum(i int) Listener {
	return no[i] // panics
}

type Listener interface {
	// GetInstance can return Empty()
	GetInstance() ident.Id
	// GetClass always returns a valid class id.
	GetClass() ident.Id
	// GetCallback() returns a valid callback id.
	GetCallback() ident.Id
	//
	GetOptions() CallbackOptions
}

type CallbackOptions int

const (
	UseTargetOnly CallbackOptions = 1 << iota
	UseAfterQueue
)

func (opt CallbackOptions) UseTargetOnly() bool {
	return opt&UseTargetOnly != 0
}
func (opt CallbackOptions) UseAfterQueue() bool {
	return opt&UseAfterQueue != 0
}

// Prototype holds properties; it supports both instance and class type things.
type Prototype interface {
	GetId() ident.Id
	//?GetType()  -> class or instance

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

type Relation interface {
	GetId() ident.Id
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
	// GetRelative panics if the property is not an object or object array.
	// it returns false if there is no relation -- a pure array or object value.
	GetRelative() (Relative, bool)
}

type Relative struct {
	Relation ident.Id
	Relates  ident.Id
	From     ident.Id
	IsRev    bool // in this future this may include a "projection" or "field name"
}

// get and set panic if the value is not of the requested type; set can return error when the value, when of the correct type, violates a property constraint
type Value interface {
	GetNum() float32
	SetNum(float32) error

	GetText() string
	SetText(string) error

	GetState() ident.Id
	SetState(ident.Id) error

	// FIX : Relations relate Objects -> instances
	GetObject() ident.Id
	SetObject(ident.Id) error
}

type Values interface {
	NumValue() int
	ValueNum(int) Value

	ClearValues() error
	AppendNum(float32) error
	AppendText(string) error
	AppendObject(ident.Id) error

	// RemoveValue(int)
}

// NOTE: ParserActions aren't id'd so, they are represented as structs.
type ParserAction struct {
	Action   ident.Id
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
