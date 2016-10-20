package standard

import (
	"github.com/ionous/mars/g"
	. "github.com/ionous/sashimi/script"
)

func init() {
	AddScript(func(s *Script) {
		s.The("actors",
			Can("attack it").And("attacking it").RequiresOne("object"),
			To("attack it", g.ReflectToTarget("report attack")))

		s.The("objects",
			Can("report attack").And("reporting attack").RequiresOne("actor"),
			To("report attack", g.Choose{
				If:   g.Equals{g.Our("player"), g.Our("actor")},
				True: g.Say("Violence isn't the answer."),
			}))

		s.Execute("attack it", Matching("attack|break|smash|hit|fight|torture {{something}}").
			Or("wreck|crack|destroy|murder|kill|punch|thump {{something}}"))
	})
}
