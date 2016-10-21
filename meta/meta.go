package meta

import "github.com/ionous/sashimi/util/ident"

type Instance Prototype
type Class Prototype

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

// FIX: this should be by id even if, in the serialized model, the action statements are in the object directly. also: callback should have its own interface i think for file and line
type Callbacks interface {
	NumCallback() int
	CallbackNum(int) Callback
}

type Event interface {
	GetId() ident.Id
	// GetEventName returns the original name given by the scripter.
	GetEventName() string
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
	GetCallback() Callback
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

	// GetParentClass returns empty id for classes if no parent;
	// panics if no class can be found for an instance.
	GetParentClass() ident.Id
	GetOriginalName() string

	NumProperty() int
	PropertyNum(int) Property

	// FindProperty by its user given name.
	FindProperty(string) (Property, bool)

	// GetProperty by the property unique id.
	GetProperty(ident.Id) (Property, bool)

	// GetPropertyByChoice evalutes all properties to find an enumeration which can store the passed choice
	GetPropertyByChoice(ident.Id) (Property, bool)
}

type Relation interface {
	GetId() ident.Id
}

type PropertyType int

const (
	InvalidProperty PropertyType = iota
	NumProperty                  // float64
	TextProperty                 // string
	StateProperty                // string.Id
	ObjectProperty               // string.Id
	ArrayProperty   = 1 << 16
)

type Generic interface{}
type Callback interface{}

type Property interface {
	GetId() ident.Id
	GetName() string
	GetType() PropertyType

	GetGeneric() Generic
	SetGeneric(Generic) error
	//GetObjectType()?
	// or maybe IsCompatible(inst) bool
	GetValue() Value
	GetValues() Values
	// GetRelative returns false if there is no relation, for example: a pure array or object value.
	GetRelative() (Relative, bool)
}

type Relative struct {
	Relation ident.Id // Relation
	Relates  ident.Id // Relates class
	From     ident.Id // From property
}

// get and set panic if the value is not of the requested type; set can return error when the value, when of the correct type, violates a property constraint
type Value interface {
	GetNum() float64
	SetNum(float64) error

	GetText() string
	SetText(string) error

	GetState() ident.Id
	SetState(ident.Id) error

	GetObject() ident.Id
	SetObject(ident.Id) error

	SetRelation(ident.Id) error
}

type Values interface {
	NumValue() int
	ValueNum(int) Value

	ClearValues() error
	AppendNum(float64) error
	AppendText(string) error
	AppendObject(ident.Id) error

	// RemoveValue(int)?
	// Slice?
}

// NOTE: ParserActions aren't id'd, so they are represented as structs.
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
