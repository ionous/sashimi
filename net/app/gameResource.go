package app

import (
	"github.com/ionous/sashimi/net/ess"
	"github.com/ionous/sashimi/net/resource"
)

// GameResource finds named sessions, or uses the "new" endpoint to create sessions via post.
// Some example uris:
// 	POST /new, create new session
// 	POST /<session>, send new input
// 	 GET /<session>/rooms/<name>/contains, list of objects
// 	 GET /<session>/classes/rooms/actions
func GameResource(sessions ess.ISessionFactory) resource.IResource {
	return resource.Wrapper{
		Finds: func(name string) (ret resource.IResource, okay bool) {
			if name == "game" {
				okay, ret = true, resource.Wrapper{
					Finds: func(name string) (ret resource.IResource, okay bool) {
						switch name {
						case "new":
							okay, ret = true, SessionCreationEndpoint(sessions)
						default:
							if res, ok := sessions.GetSession(name); ok {
								okay, ret = true, &SessionResource{res, res}
							}

						}
						return // ./game/...
					},
				}
			}
			return // .game/
		},
	}
}
