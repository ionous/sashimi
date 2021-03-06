package quips

import (
	G "github.com/ionous/sashimi/game"
)

// PlayerQuips returns the possible quips for the player to say.
// the "player" topic matches all topics; the null topic matches only null topics
func PlayerQuips(g G.Play) []G.IObject {
	qs := QuipSort{}
	con := Converse(g)
	interlocutor := con.Actor()
	if npc := interlocutor.Object(); npc.Exists() {
		qh, qm, topic := con.History(), PlayerMemory(g), con.Topic().Object()
		latest := qh.MostRecent()
		isRestrictive := latest.Exists() && latest.Is("restrictive")
		// hrmm... this is very similar to "UpdateNextQuips"
		for quips := g.Query("quips", true); quips.HasNext(); {
			quip := quips.Next()
			// Filter to quips which quip supply the interlocutor.
			if subject := quip.Get("subject").Object(); subject.Equals(npc) {
				// Filter to quips which have player comments.
				if len(quip.Get("comment").Text()) > 0 || len(quip.Get("slug").Text()) > 0 {
					// Filter quips to the current topic.
					qt := quip.Get("topic").Object()
					// the player as universal topic: applies to any topic.
					if (!qt.Exists() && !topic.Exists()) || (qt.Exists() && (topic.Equals(qt) || qt.Equals(g.The("player")))) {
						// Exclude one-time quips, checking the recollection table.
						if qm.IsQuipAllowed(quip) {
							// When following a restrictive quips, limit to those which directly follow.
							if isRestrictive && Quip(quip).Follows(latest).Directly(g) {
								qs.Add(quip)
							} else if !isRestrictive {
								// Select those which follow recent quips,
								// and those which do not have any follow constraints.
								if rank, direct := Quip(quip).Recently(qh).Follows(g); rank != 0 {
									score := qs.Add(quip)
									score.rank -= rank
									if direct {
										score.rank += 1000
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return qs.Sort()
}
