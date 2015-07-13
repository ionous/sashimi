package extensions

import (
	G "github.com/ionous/sashimi/game"
	R "github.com/ionous/sashimi/runtime"
	"github.com/ionous/sashimi/util/ident"
)

type QuipPool R.GameObjects

func GetQuipPool(g G.Play) QuipPool {
	adapt := g.(*R.GameEventAdapter)
	return QuipPool(adapt.Objects)
}

func (quips QuipPool) Interlocutor(quip *R.GameObject) (ret ident.Id) {
	if p, ok := quip.Value("Speaker").(ident.Id); ok {
		ret = p
	}
	return ret
}

func (qp QuipPool) MostRecent(qh QuipHistory) (ret *R.GameObject, okay bool) {
	id := qh.MostRecent()
	ret, okay = qp[id]
	return ret, okay
}

type visitQuips func(ident.Id, bool) bool

//
func (qp QuipPool) visitFollowers(follower ident.Id, visit visitQuips) (okay bool) {
	// search all objects
	for _, f := range qp {
		// but, only for following quips
		if isTable := f.Class().CompatibleWith("FollowingQuips"); isTable {
			// yes, this entry talks about our position relative to some other quip
			if following, ok := f.Value("Following").(ident.Id); ok && following == follower {
				// grab that other quip
				if leading, ok := f.Value("Leading").(ident.Id); ok {
					// call the visitor
					directly, _ := f.Value("DirectlyFollowing").(bool)
					if ok := visit(leading, directly); ok {
						okay = true
					}
				}
			}
		}
	}
	return okay
}

func (qp QuipPool) FollowsRecently(qh QuipHistory, follower ident.Id) (ret int) {
	isAFollower := false
	qp.visitFollowers(follower, func(leading ident.Id, directly bool) bool {
		if idx := qh.Rank(leading); (idx > ret) && (!directly || idx == QuipHistoryDepth) {
			ret = idx
		}
		isAFollower = true
		return false
	})
	if !isAFollower {
		ret = -1
	}
	return ret
}

func (qp QuipPool) FollowsDirectly(qh QuipHistory, follower ident.Id) (follows bool) {
	if mostRecent := qh.MostRecent(); !mostRecent.Empty() {
		qp.visitFollowers(follower, func(leading ident.Id, directly bool) bool {
			follows = mostRecent == leading
			return follows
		})
	}
	return follows
}

// QuipRecollects determins if the passed quip has been spoken.
func (qp QuipPool) Recollects(quip ident.Id) (recollects bool) {
	for _, r := range qp {
		if isRecollect := r.Class().CompatibleWith("Recollections"); isRecollect {
			if r.Value("Quip").(ident.Id) == quip {
				recollects = true
				break
			}
		}
	}
	return recollects
}

func VisitObjects(objects R.GameObjects, class ident.Id, visit func(*R.GameObject) bool) (okay bool) {
	for _, obj := range objects {
		if isCls := obj.Class().CompatibleWith(class); isCls {
			if visit(obj) {
				okay = true
				break
			}
		}
	}
	return okay
}

func (qp QuipPool) SpeakAfter(qh QuipHistory, newQuip *R.GameObject) (okay bool) {
	// Filter to quips which have player comments.
	if newQuip.Value("Comment").(string) != "" {
		// Exclude one-time quips, checking the recollection table.
		repeats, _ := newQuip.Value("Repeatable").(bool)
		if repeats || !qp.Recollects(newQuip.Id()) {
			// When following a restrictive quips, limit to those which directly follow.
			restricts, _ := newQuip.Value("Restrictive").(bool)
			if restricts && qp.FollowsDirectly(qh, newQuip.Id()) {
				okay = true
			} else {
				// Select those which indirect follow recent quips
				// And those which do not follow anything at all.
				if rank := qp.FollowsRecently(qh, newQuip.Id()); rank != 0 {
					okay = true
				}
			}
		}
	}
	return okay
}

// QuipList returns the possible quips for the player to say.
func (qp QuipPool) GetPlayerQuips(qh QuipHistory) (ret []*R.GameObject) {
	if lastQuip, ok := qp.MostRecent(qh); ok {
		npcId := qp.Interlocutor(lastQuip)
		VisitObjects(R.GameObjects(qp), "Quips", func(newQuip *R.GameObject) bool {
			// Filter to quips which quip supply the interlocutor.
			speaker := qp.Interlocutor(newQuip)
			if speaker == npcId {
				if qp.SpeakAfter(qh, newQuip) {
					ret = append(ret, newQuip)
				}
			}
			return false
		})
	}
	return ret
}
