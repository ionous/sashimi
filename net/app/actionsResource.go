package app

import (
	"github.com/ionous/sashimi/net/resource"
	"github.com/ionous/sashimi/runtime/api"
)

func actionResource(out resource.IBuildObjects, act api.Action) {
	nouns := act.GetNouns()
	out.NewObject(jsonId(act.GetId()), "action").
		SetAttr("act", act.GetActionName()).
		SetAttr("evt", act.GetEvent().GetEventName()).
		SetAttr("src", jsonId(nouns.Get(api.SourceNoun))).
		SetAttr("tgt", jsonId(nouns.Get(api.TargetNoun))).
		SetAttr("ctx", jsonId(nouns.Get(api.ContextNoun)))
}
