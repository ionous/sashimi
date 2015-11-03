package compiler

import (
	"fmt"
	"github.com/ionous/sashimi/compiler/call"
	E "github.com/ionous/sashimi/event"
	G "github.com/ionous/sashimi/game"
	M "github.com/ionous/sashimi/model"
	S "github.com/ionous/sashimi/source"
	"github.com/ionous/sashimi/util/errutil"
	"github.com/ionous/sashimi/util/ident"
	"io"
	"log"
	"strings"
)

type MemoryResult struct {
	*M.Model
	Calls call.MemoryStorage
}

func Compile(out io.Writer, src S.Statements) (res MemoryResult, err error) {
	calls := call.MakeMemoryStorage()
	cfg := Config{calls, out}
	if m, e := cfg.Compile(src); e != nil {
		err = e
	} else {
		res = MemoryResult{m, calls}
	}
	return
}

// Compile script statements into a "model", a form usable by the runtime.
func (cfg Config) Compile(src S.Statements) (*M.Model, error) {
	names := NewNameSource()
	rel := newRelativeFactory(names.newScope(nil))
	log := log.New(cfg.Output, "compling: ", log.Lshortfile)
	ctx := &_Compiler{
		src, names.newScope(nil),
		newClassFactory(names, rel),
		newInstanceFactory(names, log),
		rel,
		log,
		cfg.Calls,
	}
	return ctx.compile()
}

// _Compiler compiles scripts
type _Compiler struct {
	src       S.Statements
	names     NameScope
	classes   *ClassFactory
	instances *InstanceFactory
	relatives *RelativeFactory
	log       *log.Logger
	calls     call.Compiler
}

// processAssertions generates classes and instances from the assertions.
// adds to the instance and class maps
// adds to the pending instance and classes lists
func (ctx *_Compiler) processAssertions(asserts []S.AssertionStatement) (err error) {
	for err == nil && len(asserts) > 0 {
		didSomething := false
		for i := 0; i < len(asserts); i++ {
			if ok, e := ctx.processAssertBlock(asserts[i]); e != nil {
				err = errutil.Append(err, e)
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
				fields, source := pending.Fields(), pending.Source()
				e := fmt.Errorf("didn't understand how to make a `%s` called `%s`", fields.Owner, fields.Called)
				err = errutil.Append(err, SourceError(source, e))
			}
		}
	}
	return err
}

//
// Changes assert blocks into class and instance blocks.
// ( the caller keeps processing the same blocks over and over until all are resolved. )
//
func (ctx *_Compiler) processAssertBlock(assert S.AssertionStatement) (processed bool, err error) {
	fields, source := assert.Fields(), assert.Source()

	//
	longName := fields.Options["long name"]
	if longName == "" {
		longName = fields.Called
	}

	// a possible class instancing?
	// we dont mind if we dont find the owner class. it might resolve later.
	if class, ok := ctx.classes.findBySingularName(fields.Owner); ok {
		name, longName := fields.Called, longName
		if c, e := ctx.instances.addInstanceRef(class, name, longName, source); e != nil {
			err = SourceError(assert.Source(), e)
		} else if c != nil {
			processed = true
		}
	} else {
		// a possible class derivation?
		// classFactory seeds itself with the root type "kinds".
		// ultimately all classes are built as some derivation from that.
		if parent, ok := ctx.classes.findByPluralName(fields.Owner); ok {
			plural, single := longName, fields.GetOption("singular name", "")
			c, e := ctx.classes.addClassRef(parent, plural, single)
			if e != nil {
				err = SourceError(assert.Source(), e)
			} else if c != nil {
				processed = true
			}
		}
	}
	return processed || err != nil, err
}

//
func (ctx *_Compiler) compileActions(classes M.ClassMap,
) (actions M.ActionMap, events M.EventMap, err error) {
	actions, events = make(M.ActionMap), make(M.EventMap)
	//
	for _, act := range ctx.src.Actions {
		fields := act.Fields()
		if actionId, source, target, context, e := ctx.resolveAction(classes, fields); e != nil {
			err = errutil.Append(err, e)
		} else {
			// and the name of event...
			eventId, e := ctx.names.addName(fields.Event, actionId.String())
			if e != nil {
				err = errutil.Append(err, e)
				continue
			}
			// add the action; if it exists, the uniquifier should have excluded any difs, so just ignore....
			act := actions[actionId]
			if act == nil {
				if act, e = M.NewAction(actionId, fields.Action, fields.Event, source, target, context); e != nil {
					err = errutil.Append(err, e)
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

func (ctx *_Compiler) resolveAction(classes M.ClassMap, fields S.ActionAssertionFields,
) (actionId ident.Id, owner, target, context *M.ClassInfo, err error) {
	// find the primary class
	if cls, ok := classes.FindClass(fields.Source); !ok {
		e := fmt.Errorf("couldn't find class %+v", fields)
		err = errutil.Append(err, e)
	} else {
		// and the other two optional ones
		target, ok = classes[ctx.classes.singleToPlural[fields.Target]]
		if !ok && fields.Target != "" {
			e := fmt.Errorf("couldn't find class for noun %s", fields.Target)
			err = errutil.Append(err, e)
		}
		context, ok = classes[ctx.classes.singleToPlural[fields.Context]]
		if !ok && fields.Context != "" {
			e := fmt.Errorf("couldn't find class for noun %s", fields.Context)
			err = errutil.Append(err, e)
		}
		if err == nil {
			// make sure these names are unique
			owner = cls
			uniquifer := strings.Join([]string{"action", fields.Source, fields.Target, fields.Context}, "+")
			actionId, err = ctx.names.addName(fields.Action, uniquifer)
		}
	}
	return actionId, owner, target, context, err
}

//
func (ctx *_Compiler) newCallback(
	owner string,
	classes M.ClassMap,
	instances M.InstanceMap,
	action *M.ActionInfo,
	callback G.Callback,
	options M.ListenerOptions,
) (
	ret *M.ListenerCallback, err error,
) {
	if cb, e := ctx.calls.Compile(callback); e != nil {
		err = errutil.Append(e, fmt.Errorf("couldn't compile callback for `%s(%s)`", owner, action))
	} else if cls, _ := classes.FindClass(owner); cls != nil {
		ret = M.NewClassCallback(cls, action, cb.Callback, options)
	} else if inst, ok := instances.FindInstance(owner); ok {
		ret = M.NewInstanceCallback(inst, action, cb.Callback, options)
	} else {
		err = fmt.Errorf("unknown listener requested `%s(%s)`", owner, action)
	}
	return ret, err
}

//
func (ctx *_Compiler) makeActionHandlers(classes M.ClassMap, instances M.InstanceMap, actions M.ActionMap,
) (callbacks M.ActionCallbacks, err error,
) {
	for _, statement := range ctx.src.ActionHandlers {
		f := statement.Fields()
		action, ok := actions.FindActionByName(f.Action)
		if !ok {
			e := M.ActionNotFound(f.Action)
			err = errutil.Append(err, SourceError(statement.Source(), e))
			continue
		}
		var options M.ListenerOptions
		if f.Phase == E.CapturingPhase {
			options |= M.EventCapture
		} else if f.Phase == E.TargetPhase {
			options |= M.EventTargetOnly
		}

		cb, e := ctx.newCallback(f.Owner, classes, instances, action, f.Callback, options)
		if e != nil {
			err = errutil.Append(err, SourceError(statement.Source(), e))
			continue
		}
		callbacks = append(callbacks, cb)
	}
	return callbacks, err
}

//
func (ctx *_Compiler) makeEventListeners(events M.EventMap, classes M.ClassMap, instances M.InstanceMap,
) (callbacks M.ListenerCallbacks, err error,
) {
	for _, l := range ctx.src.EventHandlers {
		r := l.Fields()
		action, e := events.FindEventByName(r.Event)
		if e != nil {
			err = errutil.Append(err, SourceError(l.Source(), e))
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
		cb, e := ctx.newCallback(r.Owner, classes, instances, action, r.Callback, options)
		if e != nil {
			err = errutil.Append(err, SourceError(l.Source(), e))
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
func (ctx *_Compiler) compileAliases(instances M.InstanceMap, actions M.ActionMap) (
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
	for _, alias := range ctx.src.Aliases {
		fields := alias.Fields()
		key, phrases := fields.Key, fields.Phrases
		// alias is a noun:
		if inst, ok := instances.FindInstance(key); ok {
			id := inst.Id
			for _, name := range phrases {
				names.AddNameForId(name, id)
			}
		} else {
			// alias is an action:
			if act, ok := actions.FindActionByName(key); ok {
				// FUTURE: ensure action Parsings always involve the player object?
				// but, to know the player... might mean we couldnt run without the standard lib,
				// maybe there's user rules, or something...?
				parserAction := M.ParserAction{act, phrases}
				parserActions = append(parserActions, parserAction)
			} else {
				e := fmt.Errorf("unknown alias requested %s", key)
				err = errutil.Append(err, SourceError(alias.Source(), e))
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
func (ctx *_Compiler) compile() (*M.Model, error) {
	// create empty classes,instances,tables,actions from their assertions
	ctx.log.Println("reading assertions")
	err := ctx.processAssertions(ctx.src.Asserts)
	if err != nil {
		return nil, err
	}

	// add class primitive values;
	// queuing any we types we cant immediately resolve.
	ctx.log.Println("adding class properties")
	for _, prop := range ctx.src.Properties {
		fields := prop.Fields()
		if class, ok := ctx.classes.findByPluralName(fields.Class); !ok {
			e := SourceError(prop.Source(), ClassNotFound(fields.Class))
			err = errutil.Append(err, e)
		} else {
			if _, e := class.addProperty(prop.Source(), fields); e != nil {
				e := SourceError(prop.Source(), e)
				err = errutil.Append(err, e)
			}
		}
	}
	if err != nil {
		return nil, err
	}

	// add class enumerations
	ctx.log.Println("adding enums")
	for _, enum := range ctx.src.Enums {
		fields := enum.Fields()
		if class, ok := ctx.classes.findByPluralName(fields.Class); !ok {
			e := ClassNotFound(fields.Class)
			err = errutil.Append(err, SourceError(enum.Source(), e))
		} else if _, e := class.addEnum(fields.Name, fields.Choices, fields.Expects); e != nil {
			err = errutil.Append(err, SourceError(enum.Source(), e))
		}
	}
	if err != nil {
		return nil, err
	}

	// add class relatives
	ctx.log.Println("adding class relatives")
	for _, rel := range ctx.src.Relatives {
		fields := rel.Fields()
		if class, ok := ctx.classes.findByPluralName(fields.Class); !ok {
			err = errutil.Append(err, ClassNotFound(fields.Class))
		} else if _, e := class.addRelative(fields, rel.Source()); e != nil {
			err = errutil.Append(err, e)
		}
	}
	if err != nil {
		return nil, err
	}

	// make classes
	ctx.log.Println("making classes")
	classes, err := ctx.classes.makeClasses(ctx.relatives)
	if err != nil {
		return nil, err
	}

	// make instances
	ctx.log.Println("compiling relations")
	relations := make(M.RelationMap)
	for id, rel := range ctx.relatives.relations {
		if rel, e := rel.makeRelation(id); e != nil {
			err = errutil.Append(err, e)
		} else {
			relations[id] = rel
		}
	}
	if err != nil {
		return nil, err
	}

	multiValueData := []MultiValueData{}
	ctx.log.Println("parsing tables")
	for _, mvs := range ctx.src.MultiValues {
		fields := mvs.Fields()
		// make a table to process the rows of data for this class
		if table, e := makeValueTable(ctx.classes, fields.Owner, fields.Columns); e != nil {
			err = errutil.Append(err, SourceError(mvs.Source(), e))
		} else {
			// walk those rows
			for _, row := range fields.Rows {
				if data, e := table.addRow(ctx.instances, mvs.Source(), row); e != nil {
					err = errutil.Append(err, SourceError(mvs.Source(), e))
				} else {
					multiValueData = append(multiValueData, data)
				}
			}
		}
	}

	// make instances
	ctx.log.Println("making instances")
	partials, err := ctx.instances.makeInstances(classes, relations)
	if err != nil {
		return nil, err
	}

	// merges instance table data
	ctx.log.Println("merging instance tables")
	for _, mvd := range multiValueData {
		if e := mvd.mergeInto(partials); e != nil {
			err = errutil.Append(err, SourceError(mvd.src, e))
		}
	}
	if err != nil {
		return nil, err
	}
	// fills out the instance properties
	ctx.log.Println("setting instance properties")
	instances, tables, err := partials.makeData(ctx.src.Choices, ctx.src.KeyValues)
	if err != nil {
		return nil, err
	}

	// make actions and events
	ctx.log.Println("compiling actions")
	actions, events, err := ctx.compileActions(classes)
	if err != nil {
		return nil, err
	}

	ctx.log.Println("making action handlers")
	actHandlers, err := ctx.makeActionHandlers(classes, instances, actions)
	if err != nil {
		return nil, err
	}

	// make events listeners( in order of original declaration )
	ctx.log.Println("making event listeners")
	evtListeners, err := ctx.makeEventListeners(events, classes, instances)
	if err != nil {
		return nil, err
	}

	// create parser actions
	ctx.log.Println("compiling aliases")
	names, parserActions, err := ctx.compileAliases(instances, actions)
	if err != nil {
		return nil, err
	}

	ctx.log.Println("compiling globals")
	generators := make(M.GeneratorMap)
	for _, gen := range ctx.src.Globals {
		gf := gen.Fields()
		id := ident.MakeId(gf.Name)
		if _, exists := generators[id]; exists {
			e := fmt.Errorf("Global generator %s already exists", gf.Name)
			err = errutil.Append(err, e)
		} else {
			generators[id] = gf.Type
		}
	}

	// return the results with no error
	ctx.log.Println("compile finished")
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
		generators,
	}
	return res, nil
}
