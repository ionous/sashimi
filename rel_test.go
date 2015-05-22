package sashimi

import (
	M "github.com/ionous/sashimi/model"
	. "github.com/ionous/sashimi/script"
	x "github.com/smartystreets/goconvey/convey"
	"log"
	"os"
	"testing"
)

//
// create a single subclass called stories
func TestSimpleRelation(t *testing.T) {
	x.Convey("Test Simple Relation", t, func() {
		s := Script{}
		s.The("kinds",
			Called("gremlins"),
			Have("pets", "rocks").
				Implying("rocks", Have("o beneficent one", "gremlin")),
			// alternate, non-conflicting specification of the same relation
			HaveMany("pets", "rocks").
				// FIX? if the names don't match, this creates two views of the same relation.
				// validate the hierarchy to verify no duplicate property usage?
				Implying("rocks", HaveOne("o beneficent one", "gremlin")),
		)
		s.The("kinds", Called("rocks"), Exist())
		log := log.New(os.Stdout, "game:", log.Lshortfile)
		model, err := s.Compile(log)
		x.So(err, x.ShouldBeNil)
		model.PrintModel(t.Log)
		x.So(len(model.Relations), x.ShouldEqual, 1)
		for _, v := range model.Relations {
			x.So(v.Source(), x.ShouldEqual, "Gremlins")
			x.So(v.Destination(), x.ShouldEqual, "Rocks")
			x.So(v.Style(), x.ShouldEqual, M.OneToMany)
		}
	})
}

//
// create a single subclass called stories
func TestSimpleRelates(t *testing.T) {
	x.Convey("Test Simple Relates", t, func() {
		s := Script{}
		s.The("kinds",
			Called("gremlins"),
			Have("pets", "rocks").
				Implying("rocks", Have("o beneficent one", "gremlin")),
		)
		s.The("kinds", Called("rocks"), Exist())
		// FIX: for now the property names must match,
		// i'd prefer the signular: Has("pet", "Loofah")
		s.The("gremlin", Called("Claire"), Has("pets", "Loofah"))
		s.The("rock", Called("Loofah"), Exists())
		//
		log := log.New(os.Stdout, "game:", log.Lshortfile)
		model, err := s.Compile(log)
		x.So(err, x.ShouldBeNil)
		x.So(len(model.Instances), x.ShouldEqual, 2)

		claire, err := model.Instances.FindInstance("claire")
		x.So(err, x.ShouldBeNil)

		pets, ok := claire.ValueByName("pets")
		x.So(ok, x.ShouldBeTrue)
		petsrel := pets.(*M.RelativeValue)
		x.So(petsrel.List(), x.ShouldResemble, []string{"Loofah"})

		loofah, err := model.Instances.FindInstance("loofah")
		x.So(err, x.ShouldBeNil)

		gremlins, ok := loofah.ValueByName("o beneficent one")
		x.So(ok, x.ShouldBeTrue)
		gremlinrel := gremlins.(*M.RelativeValue)
		x.So(gremlinrel.List(), x.ShouldResemble, []string{"Claire"})

		model.PrintModel(t.Log)
	})
}
