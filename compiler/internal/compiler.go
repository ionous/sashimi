package internal

import (
	"fmt"
	"github.com/ionous/sashimi/compiler/call"
	M "github.com/ionous/sashimi/compiler/xmodel"
	E "github.com/ionous/sashimi/event"
	G "github.com/ionous/sashimi/game"
	S "github.com/ionous/sashimi/source"
	"github.com/ionous/sashimi/util/errutil"
	"github.com/ionous/sashimi/util/ident"
	"log"
	"strings"
)

// Compiler compiles scripts
type Compiler struct {
	Source    S.Statements
	Names     NameScope
	Classes   *ClassFactory
	Instances *InstanceFactory
	Relatives *RelativeFactory
	Log       *log.Logger
	Calls     call.Compiler
}

// processAssertions generates classes and instances from the assertions.
// adds to the instance and class maps
// adds to the pending instance and classes lists
func (ctx *Compiler) processAssertions(asserts []S.AssertionStatement) (err error) {
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
func (ctx *Compiler) processAssertBlock(assert S.AssertionStatement) (processed bool, err error) {
	fields, source := assert.Fields(), assert.Source()

	//
	longName := fields.Options["long name"]
	if longName == "" {
		longName = fields.Called
	}

	// a possible class instancing?
	// we dont mind if we dont find the owner class. it might resolve later.
	if class, ok := ctx.Classes.findBySingularName(fields.Owner); ok {
		name, longName := fields.Called, longName
		if c, e := ctx.Instances.addInstanceRef(class, name, longName, source); e != nil {
			err = SourceError(assert.Source(), e)
		} else if c != nil {
			processed = true
		}
	} else {
		// a possible class derivation?
		// classFactory seeds itself with the root type "kinds".
		// ultimately all classes are built as some derivation from that.
		if parent, ok := ctx.Classes.findByPluralName(fields.Owner); ok {
			plural, single := longName, fields.GetOption("singular name", "")
			c, e := ctx.Classes.addClassRef(parent, plural, single)
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
func (ctx *Compiler) compileActions(classes M.ClassMap,
) (actions M.ActionMap, events M.EventMap, err error) {
	actions, events = make(M.ActionMap), make(M.EventMap)
	//
	for _, act := range ctx.Source.Actions {
		fields := act.Fields()
		if actionId, source, target, context, e := ctx.resolveAction(classes, fields); e != nil {
			err = errutil.Append(err, e)
		} else {
			// and the name of event...
			eventId, e := ctx.Names.addName(fields.Event, actionId.String())
			if e != nil {
				err = errutil.Append(err, e)
				continue
			}
			// add the action; if it exists, the uniquifier should have excluded any difs, so just ignore....
			act := actions[actionId]
			if act == nil {
				if act, e = M.NewAction(actionId, fields.Action, eventId, source, target, context); e != nil {
					err = errutil.Append(err, e)
				} else {
					actions[actionId] = act
				}
			}
			// add the event
			if prev := events[eventId]; prev == nil {
				events[eventId] = &M.EventInfo{eventId, fields.Event, actionId}
			}
		}
	}
	return actions, events, err
}

func (ctx *Compiler) resolveAction(classes M.ClassMap, fields S.ActionAssertionFields,
) (actionId ident.Id, owner, target, context *M.ClassInfo, err error) {
	// find the primary class
	if cls, ok := classes.FindClass(fields.Source); !ok {
		e := fmt.Errorf("couldn't find class %+v", fields)
		err = errutil.Append(err, e)
	} else {
		// and the other two optional ones
		target, ok = classes[ctx.Classes.singleToPlural[fields.Target]]
		if !ok && fields.Target != "" {
			e := fmt.Errorf("couldn't find class for noun %s", fields.Target)
			err = errutil.Append(err, e)
		}
		context, ok = classes[ctx.Classes.singleToPlural[fields.Context]]
		if !ok && fields.Context != "" {
			e := fmt.Errorf("couldn't find class for noun %s", fields.Context)
			err = errutil.Append(err, e)
		}
		if err == nil {
			// make sure these names are unique
			owner = cls
			uniquifer := strings.Join([]string{"action", fields.Source, fields.Target, fields.Context}, "+")
			actionId, err = ctx.Names.addName(fields.Action, uniquifer)
		}
	}
	return actionId, owner, target, context, err
}

//
func (ctx *Compiler) newCallback(
	owner string,
	classes M.ClassMap,
	instances M.InstanceMap,
	callback G.Callback,
	options M.ListenerOptions,
) (
	ret M.ListenerCallback, err error,
) {
	if cb, e := ctx.Calls.CompileCallback(callback); e != nil {
		err = errutil.Append(e, fmt.Errorf("couldn't compile callback for `%s`", owner))
	} else if cls, _ := classes.FindClass(owner); cls != nil {
		ret = M.NewClassCallback(cls, cb, options)
	} else if inst, ok := instances.FindInstance(owner); ok {
		ret = M.NewInstanceCallback(inst, cb, options)
	} else {
		err = fmt.Errorf("unknown listener requested `%s`", owner)
	}
	return ret, err
}

//
func (ctx *Compiler) makeActionHandlers(classes M.ClassMap, instances M.InstanceMap, actions M.ActionMap,
) (callbacks M.ActionCallbacks, err error,
) {
	for _, statement := range ctx.Source.ActionHandlers {
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

		cb, e := ctx.newCallback(f.Owner, classes, instances, f.Callback, options)
		if e != nil {
			err = errutil.Append(err, SourceError(statement.Source(), e))
			continue
		}
		id := M.MakeStringId(action.ActionName)
		callbacks = append(callbacks, M.ActionCallback{id, cb})
	}
	return callbacks, err
}

//
func (ctx *Compiler) makeEventListeners(events M.EventMap, classes M.ClassMap, instances M.InstanceMap,
) (callbacks M.EventCallbacks, err error,
) {
	for _, l := range ctx.Source.EventHandlers {
		r := l.Fields()
		evt, e := events.FindEventByName(r.Event)
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
		cb, e := ctx.newCallback(r.Owner, classes, instances, r.Callback, options)
		if e != nil {
			err = errutil.Append(err, SourceError(l.Source(), e))
			continue
		}
		id := M.MakeStringId(evt.EventName)
		callbacks = append(callbacks, M.EventCallback{id, cb})
	}
	return callbacks, err
}

//
// Turn object and action aliases into noun name and parser action mappings
// FIX: instance names should go in declaration order
//
func (ctx *Compiler) compileAliases(instances M.InstanceMap, actions M.ActionMap) (
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
	for _, alias := range ctx.Source.Aliases {
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
			id := M.MakeStringId(key)
			if _, ok := actions[id]; ok {
				// FUTURE: ensure action Parsings always involve the player object?
				// but, to know the player... might mean we couldnt run without the standard lib,
				// maybe there's user rules, or something...?
				parserAction := M.ParserAction{id, phrases}
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
func (ctx *Compiler) Compile() (*M.Model, error) {
	// create empty classes,instances,tables,actions from their assertions
	ctx.Log.Println("reading assertions")
	err := ctx.processAssertions(ctx.Source.Asserts)
	if err != nil {
		return nil, err
	}

	// add class primitive values;
	// queuing any we types we cant immediately resolve.
	ctx.Log.Println("adding class properties")
	for _, prop := range ctx.Source.Properties {
		fields := prop.Fields()
		if class, ok := ctx.Classes.findByPluralName(fields.Class); !ok {
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
	ctx.Log.Println("adding enums")
	for _, enum := range ctx.Source.Enums {
		fields := enum.Fields()
		if class, ok := ctx.Classes.findByPluralName(fields.Class); !ok {
			e := ClassNotFound(fields.Class)
			err = errutil.Append(err, SourceError(enum.Source(), e))
		} else if _, e := class.addEnum(fields.Name, fields.Choices); e != nil {
			err = errutil.Append(err, SourceError(enum.Source(), e))
		}
	}
	if err != nil {
		return nil, err
	}

	// add class relatives
	ctx.Log.Println("adding class relatives")
	for _, rel := range ctx.Source.Relatives {
		fields := rel.Fields()
		if class, ok := ctx.Classes.findByPluralName(fields.Class); !ok {
			err = errutil.Append(err, ClassNotFound(fields.Class))
		} else if _, e := class.addRelative(fields, rel.Source()); e != nil {
			err = errutil.Append(err, e)
		}
	}
	if err != nil {
		return nil, err
	}

	// make classes
	ctx.Log.Println("making classes")
	classes, plurals, err := ctx.Classes.makeClasses(ctx.Relatives)
	if err != nil {
		return nil, err
	}

	// make instances
	ctx.Log.Println("compiling relations")
	relations := make(M.RelationMap)
	for id, rel := range ctx.Relatives.relations {
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
	ctx.Log.Println("parsing tables")
	for _, mvs := range ctx.Source.MultiValues {
		fields := mvs.Fields()
		// make a table to process the rows of data for this class
		if table, e := makeValueTable(ctx.Classes, fields.Owner, fields.Columns); e != nil {
			err = errutil.Append(err, SourceError(mvs.Source(), e))
		} else {
			// walk those rows
			for _, row := range fields.Rows {
				if data, e := table.addRow(ctx.Instances, mvs.Source(), row); e != nil {
					err = errutil.Append(err, SourceError(mvs.Source(), e))
				} else {
					multiValueData = append(multiValueData, data)
				}
			}
		}
	}
	if err != nil {
		return nil, err
	}

	// make instances
	ctx.Log.Println("making instances")
	partials, err := ctx.Instances.makeInstances(classes, relations)
	if err != nil {
		return nil, err
	}

	// merges instance table data
	ctx.Log.Println("merging instance tables")
	for _, mvd := range multiValueData {
		if e := mvd.mergeInto(partials); e != nil {
			err = errutil.Append(err, SourceError(mvd.src, e))
		}
	}
	if err != nil {
		return nil, err
	}
	// fills out the instance properties
	ctx.Log.Println("setting instance properties")
	instances, tables, err := partials.makeData(ctx.Source.Choices, ctx.Source.KeyValues)
	if err != nil {
		return nil, err
	}

	// make actions and events
	ctx.Log.Println("compiling actions")
	actions, events, err := ctx.compileActions(classes)
	if err != nil {
		return nil, err
	}

	ctx.Log.Println("making action handlers")
	actHandlers, err := ctx.makeActionHandlers(classes, instances, actions)
	if err != nil {
		return nil, err
	}

	// make events listeners( in order of original declaration )
	ctx.Log.Println("making event listeners")
	evtListeners, err := ctx.makeEventListeners(events, classes, instances)
	if err != nil {
		return nil, err
	}

	// create parser actions
	ctx.Log.Println("compiling aliases")
	names, parserActions, err := ctx.compileAliases(instances, actions)
	if err != nil {
		return nil, err
	}

	// FIX FIX FIX: set a generic "name" property to each instance based on "printed name"
	// really, this should combine with "long name" etc.
	//
	kinds, _ := classes.FindClass("kinds")
	directName := ident.Join(kinds.Id, ident.MakeId("name"))
	kinds.Properties[directName] = M.TextProperty{
		Id:   directName,
		Name: "name",
	}
	for _, i := range instances {
		name := i.Name
		if p, ok := i.Class.FindProperty("printed name"); ok {
			if printed, ok := i.Values[p.GetId()]; ok {
				name = printed.(string)
			}
		}
		i.Values[directName] = name
	}

	// return the results with no error
	ctx.Log.Println("compile finished")
	res := &M.Model{
		classes,
		relations,
		actions,
		events,
		parserActions,
		instances,
		actHandlers,
		evtListeners,
		tables.Tables,
		names,
		plurals,
		//generators,
	}
	return res, nil
}
