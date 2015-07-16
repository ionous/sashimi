package extensions

import (
	G "github.com/ionous/sashimi/game"
)

type QuipPool struct {
	g G.Play
}

func GetQuipPool(g G.Play) QuipPool {
	return QuipPool{g}
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

func (qp QuipPool) FollowsRecently(qh QuipHistory, follower G.IObject) (ret int) {
	isAFollower := false
	qp.visitFollowers(follower, func(leading G.IObject, directly bool) bool {
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

func (qp QuipPool) FollowsDirectly(qh QuipHistory, follower G.IObject) (follows bool) {
	if mostRecent := qh.MostRecent(qp.g); mostRecent.Exists() {
		qp.visitFollowers(follower, func(leading G.IObject, directly bool) bool {
			follows = mostRecent == leading
			return follows
		})
	}
	return follows
}

func (qp QuipPool) SpeakAfter(qh QuipHistory, newQuip G.IObject) (okay bool) {
	// Filter to quips which have player comments.
	if newQuip.Text("comment") != "" {
		// Exclude one-time quips, checking the recollection table.
		if newQuip.Is("repeatable") || !QuipMemory.Recollects(newQuip) {
			// When following a restrictive quips, limit to those which directly follow.
			if newQuip.Is("restrictive") && qp.FollowsDirectly(qh, newQuip) {
				okay = true
			} else {
				// Select those which indirect follow recent quips
				// And those which do not follow anything at all.
				if rank := qp.FollowsRecently(qh, newQuip); rank != 0 {
					okay = true
				}
			}
		}
	}
	return okay
}

// QuipList returns the possible quips for the player to say.
func (qp QuipPool) GetPlayerQuips(qh QuipHistory, npc G.IObject) (ret []G.IObject) {
	if npc.Exists() {
		qp.g.Visit("quips", func(newQuip G.IObject) bool {
			// Filter to quips which quip supply the interlocutor.
			if speaker := newQuip.Object("subject"); speaker == npc {
				if qp.SpeakAfter(qh, newQuip) {
					disallowed := IsQuipDisallowed(qp.g, newQuip)
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

// IsQuipDisallowed evaluates the quip requirements and the known facts.
// A quip requirement can allow or disallow a given quip based on whether the player knows a specific fact.
func IsQuipDisallowed(g G.Play, quip G.IObject) bool {
	disallowed := g.Visit("quip requirements",
		func(req G.IObject) (disallowed bool) {
			if req.Object("quip") == quip {
				fact, required := req.Object("fact"), req.Is("permits")
				recollects := QuipMemory.Recollects(fact)
				// the opposite behavior of required is to exclude the use of the quip
				if required != recollects {
					disallowed = true
				}
			}
			// the first disallowed quip/fact pairing stops the search because returning true stops.
			return disallowed
		})
	return disallowed
}

func GetPlayerQuips(g G.Play, qh QuipHistory, npc G.IObject) []G.IObject {
	return GetQuipPool(g).GetPlayerQuips(qh, npc)
}
