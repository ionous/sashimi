package net

// import (
// 	"encoding/json"
// 	. "github.com/ionous/sashimi/script"
// 	"github.com/ionous/sashimi/standard"
// 	"github.com/stretchr/testify/assert"
// 	"os"
// 	"testing"
// )

// func TestModelSerial(t *testing.T) {
// 	s := standard.InitStandardLibrary()
// 	s.The("story",
// 		Called("testing"),
// 		Has("author", "me"),
// 		Has("headline", "extra extra"))
// 	s.The("room",
// 		Called("the lab"),
// 		Has("description", "an empty room"))
// 	s.The("actor",
// 		Called("player"),
// 		Exists(),
// 		In("the lab"),
// 	)
// 	s.The("container",
// 		Called("cabinet"), In("the lab"),
// 		Is("openable", "closed").And("fixed in place"),
// 		Contains("glass beaker"))
// 	s.The("container",
// 		Called("the glass beaker"),
// 		Is("transparent").And("unopenable"),
// 		Has("brief", "beaker"),
// 		Contains("the eye dropper"))
// 	s.The("props",
// 		Called("droppers"),
// 		Have("drops", "num"))
// 	s.The("dropper", Called("eye dropper"), Exists(), Has("drops", 5))
// 	model, err := s.Compile(os.Stderr)
// 	if assert.NoError(t, err, "compiling test story") {
// 		ret := SerializeView(model, "Lab")
// 		b, err := json.MarshalIndent(ret, "=", "\t")
// 		if assert.NoError(t, err, "marshelling model") {
// 			os.Stderr.Write(b)
// 		}
// 	}
// }
