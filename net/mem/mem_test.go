package mem_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ionous/sashimi/_examples/stories"
	"github.com/ionous/sashimi/net"
	"github.com/ionous/sashimi/net/app"
	"github.com/ionous/sashimi/net/mem"
	"github.com/ionous/sashimi/net/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// go test -run TestMemApp
func TestMemApp(t *testing.T) {
	stories.Select("lab")

	handler := http.NewServeMux()
	handler.HandleFunc("/game/", net.HandleResource(app.GameResource(mem.NewSessions())))
	ts := httptest.NewServer(handler)
	defer ts.Close()

	g := &Helper{ts, "new"}
	if d, err := g.post(""); assert.NoError(t, err) {
		if assert.Len(t, d.Included, 0, "session starts empty") {
			if d, err := g.getMany("actors", "player", "whereabouts"); assert.NoError(t, err) && assert.Len(t, d.Included, 1, "the room") {
				if d, err := g.post("start"); assert.NoError(t, err) {
					if evts, ok := d.Data.Attributes["events"]; assert.True(t, ok, "frame has event stream") {
						require.EqualValues(t, "game", d.Data.Class)

						// read the events by re-transforming the stream
						var blocks []app.EventBlock
						b, _ := json.Marshal(evts)
						require.NoError(t, json.Unmarshal(b, &blocks))

						// test those events:
						require.True(t, len(blocks) >= 1, "reading blocks")
						require.True(t, blocks[0].Evt == "commencing", "first event should be commencing")

						// check the room
						if contents, err := g.getMany("rooms", "lab", "contents"); assert.NoError(t, err) {
							require.Len(t, contents.Data, 3, "the lab should have two objects")
							require.Len(t, contents.Included, 3, "the player should (not) be previously known, the table newly known.")
						}
						require.NoError(t, checkTable(g, 1))

						if _, err := g.post("open the glass jar"); assert.NoError(t, err) {
							require.NoError(t, checkTable(g, 1))
						}

						// take the beaker
						if _, err := g.post("take the glass jar"); assert.NoError(t, err) {
							require.NoError(t, checkTable(g, 0))
						}
						cmd := app.CommandInput{
							Action:  "show-it-to",
							Target:  "lab-assistant",
							Context: "axe"}
						if _, err := g.postCmd(cmd); assert.NoError(t, err) {

						}

						if cls, err := g.getOne("class", "droppers"); assert.NoError(t, err) {
							if parents, ok := cls.Data.Meta["classes"]; assert.True(t, ok, "has classes") {
								if parents, ok := parents.([]interface{}); assert.True(t, ok, "classes is list") {
									assert.EqualValues(t, []interface{}{"droppers", "props", "objects", "kinds"}, parents)
								}
							}
						}
					}
				}
			}
		}
	}
}

func checkTable(g *Helper, cnt int) (err error) {
	if contents, e := g.getMany("supporters", "table", "contents"); e != nil {
		err = e
	} else if len(contents.Data) != cnt {
		err = fmt.Errorf("the table should have %d objects. has %s", cnt, pretty(contents))
	}
	return err
}

type Helper struct {
	ts *httptest.Server
	id string
}

func (h *Helper) getOne(parts ...string) (doc resource.ObjectDocument, err error) {
	//"rooms", "lab", "contents"
	url := h.makeUrl(parts...)
	if resp, e := http.Get(url); e != nil {
		err = e
	} else {
		err = decodeBody(resp, &doc)
	}
	return
}

func (h *Helper) getMany(parts ...string) (doc resource.MultiDocument, err error) {
	//"rooms", "lab", "contents"
	url := h.makeUrl(parts...)
	if resp, e := http.Get(url); e != nil {
		err = e
	} else {
		err = decodeBody(resp, &doc)
	}
	return
}

func (h *Helper) post(input string) (doc resource.ObjectDocument, err error) {
	in := app.CommandInput{Input: input}
	return h.postCmd(in)
}

func (h *Helper) postCmd(in app.CommandInput) (doc resource.ObjectDocument, err error) {
	if b, e := json.Marshal(in); e != nil {
		err = e
	} else {
		postUrl := h.makeUrl()
		if resp, e := http.Post(postUrl, "application/json", bytes.NewReader(b)); e != nil {
			err = e
		} else if e := decodeBody(resp, &doc); e != nil {
			err = e
		} else {
			h.id = doc.Data.Id
		}
	}
	return doc, err
}

func (h *Helper) makeUrl(parts ...string) string {
	parts = append([]string{h.ts.URL, "game", h.id}, parts...)
	return strings.Join(parts, "/")
}

func decodeBody(resp *http.Response, d interface{}) (err error) {
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		err = json.NewDecoder(resp.Body).Decode(d)
	} else {
		body, e := ioutil.ReadAll(resp.Body)
		if e != nil {
			body = []byte(e.Error())
		}
		err = fmt.Errorf("%s %s %s", resp.Status, resp.Request.URL, body)
	}
	return err
}

func pretty(d interface{}) string {
	text, _ := json.MarshalIndent(d, "", " ")
	return string(text)
}
