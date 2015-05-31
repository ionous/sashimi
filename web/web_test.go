package web

import (
	. "github.com/ionous/sashimi/script"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"testing"
)

//
func TestNotFound(t *testing.T) {
	server := NewServer(":8088", "")
	go server.ListenAndServe()

	resp, e := http.Get("http://localhost:8088/")
	if assert.Nil(t, e) {
		defer resp.Body.Close()
		body, e := ioutil.ReadAll(resp.Body)
		if assert.Nil(t, e) {
			t.Logf("%s", body)
			assert.Equal(t, resp.StatusCode, 404)
		}
	}
}

//
func TestWebGame(t *testing.T) {
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
	server := NewServer(":8088", "")
	match := regexp.MustCompile("text/([^/]+)/")
	go server.ListenAndServe()

	resp, e := http.PostForm("http://localhost:8088/text/new", url.Values{"q": {""}})
	if assert.NoError(t, e) {
		defer resp.Body.Close()
		if body, e := ioutil.ReadAll(resp.Body); assert.NoError(t, e) {
			if assert.Equal(t, resp.StatusCode, 200) {
				// we expect that we've been redirected
				got := match.FindStringSubmatch(resp.Request.URL.Path)
				sess := got[1]
				if len(sess) < 16 {
					t.Fatal(got)
				}
				t.Logf("Received '%s':%s", sess, string(body))
			}
		}
	}
}
