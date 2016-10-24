package internal

import (
	"github.com/ionous/mars/rt"
	"github.com/ionous/mars/scope"
	E "github.com/ionous/sashimi/event"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/errutil"
	"github.com/ionous/sashimi/util/ident"
	"github.com/ionous/sashimi/util/sbuf"
)

type Dispatch struct {
	game   *Game
	nouns  meta.Nouns
	values []meta.Generic
	after  []meta.Callback
}

func (d *Dispatch) RunActionNow(cb meta.Callback, hint ident.Id) (err error) {
	// MARS fix? kind of annoying we have to create a scope just to inject the hint
	// possibly listener scope should allow SetHint
	newScope := scope.Make(d.game.Rtm, scope.NewListener(
		d.game.Model, d.nouns, hint, d.values))
	if call, ok := cb.(rt.Execute); !ok {
		err = errutil.New("callback not of execute type", sbuf.Type{cb})
	} else {
		err = call.Execute(newScope)
	}
	return
}

// queue for running after the default actions
func (d *Dispatch) RunActionLater(cb meta.Callback, _ ident.Id) {
	d.after = append(d.after, cb)
}

func (d *Dispatch) GetClass(cls ident.Id) (ret meta.Class, okay bool) {
	return d.game.Model.GetClass(cls)
}

func (d *Dispatch) GetParent(obj meta.Instance) (ret meta.Instance, okay bool) {
	if next, _, haveParent := d.game.LookupParent(obj); haveParent {
		okay, ret = true, next
	}
	return
}

// target: class or instance id
// note: we get multiple dispatch calls for each event: capture, target, and bubble.
func (d *Dispatch) DispatchEvent(evt E.IEvent, target ident.Id) (err error) {
	if src, ok := d.game.Model.GetEvent(evt.Id()); ok {
		if ls, ok := src.GetListeners(true); ok {
			err = E.Capture(evt, NewGameListeners(d, evt, target, ls))
		}
		if err == nil {
			if ls, ok := src.GetListeners(false); ok {
				err = E.Bubble(evt, NewGameListeners(d, evt, target, ls))
			}
		}
	}
	return
}

// run "after" actions queued by RunActionLater
func (d *Dispatch) RunAfterActions(run rt.Runtime) (err error) {
	if after := d.after; len(after) > 0 {
		for _, cb := range after {
			if call, ok := cb.(rt.Execute); !ok {
				err = errutil.New("Game.Run, after callback not of execute type", sbuf.Type{cb})
				break
			} else if e := call.Execute(run); e != nil {
				err = e
				break
			}
		}
	}
	return
}
