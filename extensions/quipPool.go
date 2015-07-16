package extensions

import (
	G "github.com/ionous/sashimi/game"
)

type QuipPool struct {
	g G.Play
	*Conversation
}

func GetQuipPool(g G.Play) QuipPool {
	return QuipPool{g, g.Global("conversation").(*Conversation)}
}

type visitQuips func(G.IObject, bool) bool

//
func (qp QuipPool) visitFollowers(follower G.IObject, visit visitQuips) bool {
	// search all following quips
	return qp.g.Visit("FollowingQuips", func(obj G.IObject) (okay bool) {
		// yes, this entry talks about our position relative to some other quip
		if following := obj.Object("Following"); following.Exists() && following == follower {
			// grab that other quip
			if leading := obj.Object("Leading"); leading.Exists() {
				// call the visitor
				directly := obj.Is("DirectlyFollowing")
				if ok := visit(leading, directly); ok {
					okay = true
				}
			}
		}
		return okay
	})
}

func (qp QuipPool) FollowsRecently(follower G.IObject) (ret int) {
	isAFollower := false
	qp.visitFollowers(follower, func(leading G.IObject, directly bool) bool {
		if idx := qp.History.Rank(leading); (idx > ret) && (!directly || idx == QuipHistoryDepth) {
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

func (qp QuipPool) FollowsDirectly(follower G.IObject) (follows bool) {
	if mostRecent := qp.History.MostRecent(qp.g); mostRecent.Exists() {
		qp.visitFollowers(follower, func(leading G.IObject, directly bool) bool {
			follows = mostRecent == leading
			return follows
		})
	}
	return follows
}

func (qp QuipPool) SpeakAfter(newQuip G.IObject) (okay bool) {
	// Filter to quips which have player comments.
	if newQuip.Text("comment") != "" {
		// Exclude one-time quips, checking the recollection table.
		if newQuip.Is("repeatable") || !qp.Memory.Recollects(newQuip) {
			// When following a restrictive quips, limit to those which directly follow.
			if newQuip.Is("restrictive") && qp.FollowsDirectly(newQuip) {
				okay = true
			} else {
				// Select those which indirect follow recent quips
				// And those which do not follow anything at all.
				if rank := qp.FollowsRecently(newQuip); rank != 0 {
					okay = true
				}
			}
		}
	}
	return okay
}

// QuipList returns the possible quips for the player to say.
func (qp QuipPool) GetPlayerQuips() (ret []G.IObject) {
	if npc, ok := qp.Interlocutor.Get(); ok {
		qp.g.Visit("quips", func(newQuip G.IObject) bool {
			speaker := newQuip.Object("subject")
			// Filter to quips which quip supply the interlocutor.
			if speaker == npc {
				after := qp.SpeakAfter(newQuip)
				if after {
					disallowed := qp.Memory.IsQuipDisallowed(qp.g, newQuip)
					if !disallowed {
						ret = append(ret, newQuip)
					}
				}
			}
			return false
		})
	}
	return ret
}

func GetPlayerQuips(g G.Play) []G.IObject {
	qp := GetQuipPool(g)
	return qp.GetPlayerQuips()
}
