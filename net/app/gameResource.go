package app

import (
	"github.com/ionous/sashimi/net/resource"

	"github.com/ionous/sashimi/net/ess"
)

// GameResource finds named sessions, or uses the "new" endpoint to create sessions via post.
// Some example uris:
// 	POST /new, create new session
// 	POST /<session>, send new input
// 	 GET /<session>/rooms/<name>/contains, list of objects
// 	 GET /<session>/classes/rooms/actions
func GameResource(sessions ess.ISessionResourceFactory) resource.IResource {
	return resource.Wrapper{
		Finds: func(name string) (ret resource.IResource, okay bool) {
			if name == "game" {
				okay, ret = true, resource.Wrapper{
					Finds: func(name string) (ret resource.IResource, okay bool) {
						switch name {
						case "new":
							okay, ret = true, NewSessionResource(sessions)
						default:
							ret, okay = sessions.GetSession(name)
						}
						return // ./game/...
					},
				}
			}
			return // .game/
		},
	}
}
