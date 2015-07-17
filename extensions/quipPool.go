package extensions

import (
	G "github.com/ionous/sashimi/game"
)

type followsCb func(leads G.IObject, directly bool) bool

// evaluate all quips which constrain this clip
// ex. QuipHelp(x).DirectlyFollows(y) -> visit(x) will call cb(y, true)
func visitFollowConstraints(g G.Play, follower G.IObject, cb followsCb) bool {
	// search all following quips
	return g.Visit("FollowingQuips", func(obj G.IObject) (okay bool) {
		// yes, this entry talks about our position relative to some other quip
		if following := obj.Object("following"); following.Exists() && following == follower {
			// grab that other quip
			if leading := obj.Object("leading"); leading.Exists() {
				// call the visitor
				directly := obj.Is("directly following")
				if ok := cb(leading, directly); ok {
					okay = true
				}
			}
		}
		return okay
	})
}

// QuipHelp provides object oriented functions for evaluating quip relations
type QuipHelp struct {
	quip G.IObject
}

// Quip returns QuipHelp
func Quip(quip G.IObject) QuipHelp {
	return QuipHelp{quip}
}

// Follows provides information about the order of quips in a conversation.
func (q QuipHelp) Follows(leader G.IObject) DirectInfo {
	return DirectInfo{q.quip, leader}
}

// Directly returns true if the Quip should only be displayed after Follows
func (info DirectInfo) Directly(g G.Play) bool {
	return visitFollowConstraints(g, info.follower, func(leading G.IObject, directly bool) bool {
		return directly && info.leader == leading
	})
}

// Recently provides information about the order of quips in a conversation.
func (q QuipHelp) Recently(history QuipHistory) RecentInfo {
	return RecentInfo{q.quip, history}
}

// Follows ranks the Quip against all recent history.
// Returns -1 if the passed quip follows no other quip;
// returns 0 if the passed quip follows something, but not one of the recent quips;
// otherwise, the higher the number, the more recent the quip that it follows.
func (info RecentInfo) Follows(g G.Play) (ret int, direct bool) {
	isAFollower := false
	visitFollowConstraints(g, info.quip, func(leading G.IObject, directly bool) bool {
		// find the most recent (highest rank) quip.
		// we only want to consider directly following quips if we are indeed directly following them.
		if rank := info.qh.Rank(leading); (rank > ret) && (!directly || rank == QuipHistoryDepth) {
			ret, direct = rank, directly
		}
		isAFollower = true
		return false // searches all
	})
	if !isAFollower {
		ret = -1
	}
	return ret, direct
}

// DirectInfo provides object oriented functions for evaluating quip relations
type DirectInfo struct {
	follower, leader G.IObject
}

// RecentInfo provides object oriented functions for evaluating quip relations
type RecentInfo struct {
	quip G.IObject
	qh   QuipHistory
}
