package app

import (
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/net/resource"
)

func actionResource(out resource.IBuildObjects, act meta.Action) {
	nouns := act.GetNouns()
	out.NewObject(jsonId(act.GetId()), "action").
		SetAttr("act", act.GetActionName()).
		SetAttr("evt", act.GetEvent().GetEventName()).
		SetAttr("src", jsonId(nouns.Get(meta.SourceNoun))).
		SetAttr("tgt", jsonId(nouns.Get(meta.TargetNoun))).
		SetAttr("ctx", jsonId(nouns.Get(meta.ContextNoun)))
}
