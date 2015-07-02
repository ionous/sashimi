package compiler

import (
	"fmt"
	E "github.com/ionous/sashimi/event"
	G "github.com/ionous/sashimi/game"
	M "github.com/ionous/sashimi/model"
	S "github.com/ionous/sashimi/source"
	"github.com/ionous/sashimi/util/errutil"
	"io"
	"log"
	"strings"
)

type ErrorLog struct {
	*log.Logger
}

// before appending the error in the normal way, this.log.
func (this ErrorLog) AppendError(err error, e error) error {
	this.Output(2, e.Error())
	return errutil.Append(err, e)
}

//
// create a new instance of script results
func Compile(out io.Writer, src S.Blocks) (*M.Model, error) {
	names := NewNameSource()
	ctx := &Context{
		src, names.newScope(nil),
		newClassFactory(names),
		newInstanceFactory(names),
		newRelativeFactory(names.newScope(nil)),
		&ErrorLog{log.New(out, "compling: ", log.Lshortfile)},
	}
	return ctx.compile()
}

//
// internal for compiling a script
type Context struct {
	src       S.Blocks
	names     NameScope
	classes   *ClassFactory
	instances *InstanceFactory
	relatives *RelativeFactory
	log       *ErrorLog
}

//
// generate classes and instances from the assertions
// adds to the instance and class maps
// adds to the pending instance and classes lists
func (this *Context) processAssertions(asserts []S.AssertionStatement) (err error) {
	for err == nil && len(asserts) > 0 {
		didSomething := false
		for i := 0; i < len(asserts); i++ {
			if ok, e := this.processAssertBlock(asserts[i]); e != nil {
				err = this.log.AppendError(err, e)
			} else if ok {
				last := len(asserts) - 1
				if last != i {
					asserts[i] = asserts[last]
				}
				asserts = asserts[:last]
				didSomething = true
				continue
			}
		}
		// went through a full loop without being able to process anything
		if !didSomething {
			for _, pending := range asserts {
				e := fmt.Errorf("couldn't resolve the assertion for `%v`", pending)
				err = this.log.AppendError(err, e)
			}
		}
	}
	return err
}

//
// Changes assert blocks into class and instance blocks.
// ( the caller keeps processing the same blocks over and over until all are resolved. )
//
func (this *Context) processAssertBlock(assert S.AssertionStatement) (processed bool, err error) {
	owner := assert.Owner()

	// a possible class instancing?
	// we dont mind if we dont find the owner class. it might resolve later.
	if class, ok := this.classes.findBySingularName(owner); ok {
		if c, e := this.instances.addInstanceRef(assert.ShortName(), class.id, assert.Options(), assert.Source()); e != nil {
			err = e
		} else if c != nil {
			processed = true
		}
	} else {
		// a possible class derivation?
		// classFactory seeds itself with the root type "kinds".
		// ultimately all classes are built as some derivation from that.
		if parent, ok := this.classes.findByPluralName(owner); ok {
			c, e := this.classes.addClassRef(parent, assert.FullName(), assert.Options())
			if e != nil {
				err = e
			} else if c != nil {
				processed = true
			}
		}
	}
	return processed || err != nil, err
}

//
func (this *Context) compileActions(classes M.ClassMap,
) (actions M.ActionMap, events M.EventMap, err error) {
	actions, events = make(M.ActionMap), make(M.EventMap)
	//
	for _, act := range this.src.Actions {
		fields := act.Fields()
		if actionId, source, target, context, e := this.resolveAction(classes, fields); e != nil {
			err = this.log.AppendError(err, e)
		} else {
			// and the name of event...
			eventId, e := this.names.addName(fields.Event, actionId.String())
			if e != nil {
				err = this.log.AppendError(err, e)
				continue
			}
			// add the action; if it exists, the uniquifier should have excluded any difs, so just ignore....
			act := actions[actionId]
			if act == nil {
				if act, e = M.NewAction(actionId, fields.Action, fields.Event, source, target, context); e != nil {
					err = this.log.AppendError(err, e)
				} else {
					actions[actionId] = act
				}
			}
			// add the event
			if prev := events[eventId]; prev == nil {
				events[eventId] = act
			}
		}
	}
	return actions, events, err
}

func (this *Context) resolveAction(classes M.ClassMap, fields S.ActionAssertionFields,
) (actionId M.StringId, owner, target, context *M.ClassInfo, err error) {
	// find the primary class
	if cls, ok := classes.FindClass(fields.Source); !ok {
		e := fmt.Errorf("couldn't find class %+v", fields)
		err = this.log.AppendError(err, e)
	} else {
		// and the other two optional ones
		target, ok = classes[this.classes.singleToPlural[fields.Target]]
		if !ok && fields.Target != "" {
			e := fmt.Errorf("couldn't find class for noun %s", fields.Target)
			err = this.log.AppendError(err, e)
		}
		context, ok = classes[this.classes.singleToPlural[fields.Context]]
		if !ok && fields.Context != "" {
			e := fmt.Errorf("couldn't find class for noun %s", fields.Context)
			err = this.log.AppendError(err, e)
		}
		if err == nil {
			// make sure these names are unique
			owner = cls
			uniquifer := strings.Join([]string{"action", fields.Source, fields.Target, fields.Context}, "+")
			actionId, err = this.names.addName(fields.Action, uniquifer)
		}
	}
	return actionId, owner, target, context, err
}

//
func (this *Context) newCallback(
	owner string,
	classes M.ClassMap,
	instances M.InstanceMap,
	action *M.ActionInfo,
	cb G.Callback,
	options M.ListenerOptions,

) (ret *M.ListenerCallback, err error,
) {
	if cls, _ := classes.FindClass(owner); cls != nil {
		ret = M.NewClassCallback(cls, action, cb, options)
	} else if inst, ok := instances.FindInstance(owner); ok {
		ret = M.NewInstanceCallback(inst, action, cb, options)
	} else {
		err = fmt.Errorf("unknown listener requested `%s(%s)`", owner, action)
	}
	return ret, err
}

//
func (this *Context) makeActionHandlers(classes M.ClassMap, instances M.InstanceMap, actions M.ActionMap,
) (callbacks M.ActionCallbacks, err error,
) {
	for _, statement := range this.src.ActionHandlers {
		f := statement.Fields()
		action, ok := actions.FindActionByName(f.Action)
		if !ok {
			err = this.log.AppendError(err, M.ActionNotFound(f.Action))
			continue
		}
		var options M.ListenerOptions
		if f.Phase == E.CapturingPhase {
			options |= M.EventCapture
		} else if f.Phase == E.TargetPhase {
			options |= M.EventTargetOnly
		}

		cb, e := this.newCallback(f.Owner, classes, instances, action, f.Callback, options)
		if e != nil {
			err = this.log.AppendError(err, e)
			continue
		}
		callbacks = append(callbacks, cb)
	}
	return callbacks, err
}

//
func (this *Context) makeEventListeners(events M.EventMap, classes M.ClassMap, instances M.InstanceMap,
) (callbacks M.ListenerCallbacks, err error,
) {
	for _, l := range this.src.EventHandlers {
		r := l.Fields()
		action, e := events.FindEventByName(r.Event)
		if e != nil {
			err = this.log.AppendError(err, e)
			continue
		}
		var options M.ListenerOptions
		if r.Captures() {
			options |= M.EventCapture
		}
		if r.OnlyTargets() {
			options |= M.EventTargetOnly
		}
		if r.RunsAfter() {
			options |= M.EventQueueAfter
		}
		cb, e := this.newCallback(r.Owner, classes, instances, action, r.Callback, options)
		if e != nil {
			err = this.log.AppendError(err, e)
			continue
		}
		callbacks = append(callbacks, cb)
	}
	return callbacks, err
}

//
// Turn object and action aliases into noun name and parser action mappings
// FIX: instance names should go in declaration order
//
func (this *Context) compileAliases(instances M.InstanceMap, actions M.ActionMap) (
	names M.NounNames,
	parserActions []M.ParserAction,
	err error,
) {
	names = make(M.NounNames)

	// first: add the full names of each instance at highest ranks
	for k, _ := range instances {
		parts := k.Split()
		fullName := strings.Join(parts, " ")
		names.AddNameForId(fullName, k)
	}

	// then: add all "is known as"
	for _, alias := range this.src.Aliases {
		key, phrases := alias.Key(), alias.Phrases()
		if inst, ok := instances.FindInstance(key); ok {
			id := inst.Id()
			for _, name := range phrases {
				names.AddNameForId(name, id)
			}
		} else {
			if act, ok := actions.FindActionByName(key); ok {
				// FUTURE: ensure action Parsings always involve the player object?
				// but, to know the player... might mean we couldnt run without the standard lib,
				// maybe there's user rules, or something...?
				parserAction := M.NewParserAction(act, phrases)
				parserActions = append(parserActions, parserAction)
			} else {
				e := fmt.Errorf("unknown alias requested %s", key)
				err = this.log.AppendError(err, e)
			}
		}
	}

	// finally: add the parts as lesser ranks
	for k, _ := range instances {
		parts := k.Split()
		if len(parts) > 0 {
			for _, p := range parts {
				names.AddNameForId(p, k)
			}
		}
	}

	return names, parserActions, err
}

//
func (this *Context) compile() (*M.Model, error) {
	// create empty classes,instances,tables,actions from their assertions
	this.log.Println("reading assertions")
	err := this.processAssertions(this.src.Asserts)
	if err != nil {
		return nil, err
	}

	// add class primitive values;
	// queuing any we types we cant immediately resolve.
	pendingPointers := []S.PropertyStatement{}
	this.log.Println("adding class properties")
	for _, prop := range this.src.Properties {
		fields := prop.Fields()
		if class, ok := this.classes.findByPluralName(fields.Class); !ok {
			err = this.log.AppendError(err, ClassNotFound(fields.Class))
		} else {
			if prim, e := class.addPrimitive(prop.Source(), fields); e != nil {
				err = this.log.AppendError(err, e)
			} else if prim == nil {
				pendingPointers = append(pendingPointers, prop)
			}
		}
	}
	if err != nil {
		return nil, err
	}

	// add class enumerations
	this.log.Println("adding enums")
	for _, enum := range this.src.Enums {
		fields := enum.Fields()
		if class, ok := this.classes.findByPluralName(fields.Class); !ok {
			err = this.log.AppendError(err, ClassNotFound(fields.Class))
		} else {
			if _, e := class.addEnum(fields.Name, fields.Choices, fields.Expects); e != nil {
				err = this.log.AppendError(err, e)
			}
		}
	}
	if err != nil {
		return nil, err
	}

	// add class primitive values
	this.log.Println("adding class relatives")
	for _, rel := range this.src.Relatives {
		fields := rel.Fields()
		if class, ok := this.classes.findByPluralName(fields.Class); !ok {
			err = this.log.AppendError(err, ClassNotFound(fields.Class))
		} else {
			if e := class.addRelative(fields, rel.Source()); e != nil {
				err = this.log.AppendError(err, e)
			}
		}
	}
	if err != nil {
		return nil, err
	}

	// add foreign keys
	this.log.Println("adding class pointers")
	for _, prop := range pendingPointers {
		fields := prop.Fields()
		if class, ok := this.classes.findByPluralName(fields.Class); !ok {
			err = this.log.AppendError(err, ClassNotFound(fields.Class))
		} else if e := class.addPointer(fields, prop.Source()); e != nil {
			err = this.log.AppendError(err, e)
		}
	}

	// make classes
	this.log.Println("making classes")
	classes, err := this.classes.makeClasses(this.relatives)
	if err != nil {
		return nil, err
	}

	// make instances
	this.log.Println("compiling relations")
	relations := make(M.RelationMap)
	for id, rel := range this.relatives.relations {
		if rel, e := rel.makeRelation(id); e != nil {
			err = this.log.AppendError(err, e)
		} else {
			relations[id] = rel
		}
	}
	if err != nil {
		return nil, err
	}

	multiValueData := []MultiValueData{}
	this.log.Println("parsing tables")
	for _, mvs := range this.src.MultiValues {
		fields := mvs.Fields()
		// have to delay instance data until after the instnaces have been created.
		if class, ok := classes.FindClassBySingular(fields.Owner); ok {
			// make a table to process the rows of data for this class
			if table, e := MakeValueTable(class, fields.Columns); e != nil {
				err = this.log.AppendError(err, e)
			} else {
				// walk those rows
				for _, row := range fields.Rows {
					if data, e := table.AddRow(this.instances, mvs.Source(), row); e != nil {
						err = this.log.AppendError(err, e)
					} else {
						multiValueData = append(multiValueData, data)
					}
				}
			}
		} else {
			e := mvs.Source().Errorf("couldn't find a class or instance for %s", fields.Owner)
			err = this.log.AppendError(err, e)
		}
	}

	// make instances
	this.log.Println("making instances")
	partials, err := this.instances.makeInstances(this.log, classes, relations)
	if err != nil {
		return nil, err
	}

	// fills out the instance properties
	this.log.Println("setting instance properties")
	instances, tables, err := partials.makeData(this.src.Choices, this.src.KeyValues)
	if err != nil {
		return nil, err
	}

	// merges instance table data
	this.log.Println("merging instance tables")
	for _, mvd := range multiValueData {
		if e := mvd.mergeInto(instances); e != nil {
			err = errutil.Append(err, e)
		}
	}
	if err != nil {
		return nil, err
	}

	// make actions and events
	this.log.Println("compiling actions")
	actions, events, err := this.compileActions(classes)
	if err != nil {
		return nil, err
	}

	this.log.Println("making action handlers")
	actHandlers, err := this.makeActionHandlers(classes, instances, actions)
	if err != nil {
		return nil, err
	}

	// make events listeners( in order of original declaration )
	this.log.Println("making event listeners")
	evtListeners, err := this.makeEventListeners(events, classes, instances)
	if err != nil {
		return nil, err
	}

	// create parser actions
	this.log.Println("compiling aliases")
	names, parserActions, err := this.compileAliases(instances, actions)
	if err != nil {
		return nil, err
	}

	// return the results with no error
	this.log.Println("compile finished")
	res := &M.Model{
		classes,
		relations,
		actions,
		events,
		parserActions,
		instances,
		names,
		actHandlers,
		evtListeners,
		tables,
	}
	return res, nil
}
