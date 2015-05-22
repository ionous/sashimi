package model

import (
	x "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestStringIds(t *testing.T) {
	x.Convey("Given the script", t, func() {
		x.So(MakeStringId("apples"), x.ShouldEqual, "Apples")
		x.So(MakeStringId("Apples"), x.ShouldEqual, "Apples")
		x.So(MakeStringId("apple turnover"), x.ShouldEqual, "AppleTurnover")
		x.So(MakeStringId("Apple Turnover"), x.ShouldEqual, "AppleTurnover")
		x.So(MakeStringId("Apple turnover"), x.ShouldEqual, "AppleTurnover")
		x.So(MakeStringId("APPLE TURNOVER"), x.ShouldEqual, "AppleTurnover")
		x.So(MakeStringId("PascalCase"), x.ShouldEqual, "PascalCase")
		x.So(MakeStringId("camelCase"), x.ShouldEqual, "CamelCase")
		x.So(MakeStringId("the article"), x.ShouldEqual, "Article")
		x.So(MakeStringId("something-like-this"), x.ShouldEqual, "SomethingLikeThis")
		x.So(MakeStringId("ALLCAPS"), x.ShouldEqual, "Allcaps")
		x.So(MakeStringId("whaTAboutThis"), x.ShouldEqual, "WhaTaboutThis")
		x.So(MakeStringId("AppleTurnover").Split(), x.ShouldResemble, []string{"apple", "turnover"})
		x.So(MakeStringId("Title").Split(), x.ShouldResemble, []string{"title"})
	})
}
