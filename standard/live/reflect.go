package live

import (
	G "github.com/ionous/sashimi/game"
)

// reflect to the passed action passing the actors's current whereabouts.
// will have to become more sophisticated for being inside a box.
func ReflectToLocation(g G.Play, action string) {
	actor := g.The("actor")
	target := actor.Object("whereabouts")
	//g.Log("reflecting", action, actor, target)
	target.Go(action, actor)
}

// ReflectToTarget runs the passed action, flipping the source and target.
func ReflectToTarget(g G.Play, action string) {
	source := g.The("action.Source")
	target := g.The("action.Target")
	//g.Log("reflecting", action, source, target)
	target.Go(action, source)
}

// ReflectWithContext runs the passed action, shifting to target, context, source.
// FIX: i think it'd be better to first use ReflectToTarget, keeping the context as the third parameter
// and then ReflectToContext, possibly re-swapping source and target.
func ReflectWithContext(g G.Play, action string) {
	source := g.The("action.Source")
	target := g.The("action.Target")
	context := g.The("action.Context")
	//g.Log("reflecting", action, source, target, context)
	target.Go(action, context, source)
}
