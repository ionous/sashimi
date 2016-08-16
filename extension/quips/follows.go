package quips

import (
	G "github.com/ionous/sashimi/game"
)

type followsCb func(leads G.IObject, directly bool) bool

// evaluate all quips which constrain this clip
// ex. for QuipHelp(x).DirectlyFollows(y), then visit(x) will call cb(y, true)
func visitFollowConstraints(g G.Play, follower G.IObject, cb followsCb) (okay bool) {
	// search all following quips
	for quips := g.Query("following quips", true); quips.HasNext(); {
		quip := quips.Next()
		// yes, this entry talks about our position relative to some other quip
		if following := quip.Get("following").Object(); following.Exists() && following.Equals(follower) {
			// grab that other quip
			if leading := quip.Get("leading").Object(); leading.Exists() {
				// call the visitor
				directly := quip.Is("directly following")
				if ok := cb(leading, directly); ok {
					okay = true
					break
				}
			}
		}
	}
	return
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

// Directly returns true if the quip should only be displayed after Follows
func (info DirectInfo) Directly(g G.Play) bool {
	return visitFollowConstraints(g, info.follower, func(leading G.IObject, directly bool) bool {
		return directly && info.leader.Equals(leading)
	})
}

// Recently provides information about the order of quips in a conversation.
func (q QuipHelp) Recently(history QuipHistory) RecentInfo {
	return RecentInfo{q.quip, history}
}

// Follows ranks the Quip against all recent history.
// Returns -1 if the quip follows no other quip;
// Returns 0 if the quip follows something, but not one of the recent quips;
// Otherwise, the higher the number, the more recent the quip that it follows.
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
