package net

import (
	"github.com/ionous/sashimi/net/simple"
	. "github.com/ionous/sashimi/script"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
)

//
func TestSimpleStartup(t *testing.T) {
	ts := httptest.NewServer(simple.NewSimpleServer())
	defer ts.Close()

	resp, e := http.Get(ts.URL)
	if assert.Nil(t, e) {
		defer resp.Body.Close()
		body, e := ioutil.ReadAll(resp.Body)
		if assert.Nil(t, e) {
			t.Logf("%s", body)
			assert.Equal(t, 200, resp.StatusCode)
		}
	}
}

//
func TestSimpleGame(t *testing.T) {
	// we need a story to play
	AddScript(func(s *Script) {
		s.The("story",
			Called("testing"),
			Has("author", "me"),
			Has("headline", "extra extra"))
		s.The("room",
			Called("somewhere"),
			Has("description", "an empty room"),
		)
	})
	//
	ts := httptest.NewServer(simple.NewSimpleServer())
	defer ts.Close()
	URL := ts.URL + "/game/new"
	match := regexp.MustCompile("^/game/([^/]+)$")
	resp, e := http.PostForm(URL, url.Values{"q": {""}})
	if assert.NoError(t, e) {
		defer resp.Body.Close()
		if body, e := ioutil.ReadAll(resp.Body); assert.NoError(t, e) {
			if assert.Equal(t, 200, resp.StatusCode) {
				// we expect that we've been redirected
				got := match.FindStringSubmatch(resp.Request.URL.Path)
				t.Log("Got", got, "from:", resp.Request.URL.Path)
				if assert.Len(t, got, 2) {
					sess := got[1]
					if len(sess) < 16 {
						t.Fatal(got)
					}
					t.Logf("Received '%s':%s", sess, string(body))
				}
			}
		}
	}
}
