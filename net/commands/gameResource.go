package commands

import (
	"github.com/ionous/sashimi/net/resource"
	"github.com/ionous/sashimi/net/session"
)

//
// Finds named sessions, or uses the "new" endpoint to create sessions von post.
//
func GameResource(sessions *session.Sessions) resource.IResource {
	return resource.Wrapper{
		Finds: func(name string) (ret resource.IResource, okay bool) {
			if name == "game" {
				okay, ret = true, resource.Wrapper{
					Finds: func(name string) (ret resource.IResource, okay bool) {
						switch name {
						case "new":
							okay, ret = true, NewSessionResource(sessions)
						default:
							if sd, ok := sessions.Session(name); ok {
								session := sd.(*CommandSession)
								okay, ret = true, &SessionResource{session, session}
							}
						}
						return ret, okay
					},
				}
			}
			return ret, okay
		},
	}
}
