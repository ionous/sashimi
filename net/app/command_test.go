package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ionous/sashimi/_examples/stories"
	"github.com/ionous/sashimi/compiler/call"
	"github.com/ionous/sashimi/net/resource"
	"github.com/ionous/sashimi/net/session"
	"github.com/ionous/sashimi/script"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

//
func TestNetApp(t *testing.T) {
	stories.Select("lab")

	calls := call.MakeMemoryStorage()
	sessions := session.NewSessions(
		func(id string) (ret session.SessionData, err error) {
			// FIX: it's very silly to have to init and compile each time.
			// the reason is because relations change the original model.
			if m, e := script.InitScripts().CompileCalls(ioutil.Discard, calls); e != nil {
				err = e
			} else {
				ret, err = NewCommandSession(id, m, calls)
			}
			return ret, err
		})

	handler := http.NewServeMux()
	handler.HandleFunc("/game/", NewGameHandler(sessions))
	ts := httptest.NewServer(handler)
	defer ts.Close()

	g := &Helper{ts, "new"}
	if d, err := g.post(""); assert.NoError(t, err) {
		if assert.Len(t, d.Included, 0, "the player and the room") {
			d, err := g.post("start")
			if assert.NoError(t, err) && assert.Len(t, d.Included, 2, "the player and the room") {

				if _, ok := d.Data.Attributes["events"]; assert.True(t, ok, "frame has event stream") {

					if !assert.EqualValues(t, "game", d.Data.Class) {
						return
					}
					// check the room
					if contents, err := g.getMany("rooms", "lab", "contents"); assert.NoError(t, err) {
						assert.Len(t, contents.Data, 3, "the lab should have two objects")
						assert.Len(t, contents.Included, 2, "the player should be previously known, the table newly known.")
					}
					assert.NoError(t, checkTable(g, 1))

					if _, err := g.post("open the glass jar"); assert.NoError(t, err) {
						require.NoError(t, checkTable(g, 1))
					}

					// take the beaker
					if _, err := g.post("take the glass jar"); assert.NoError(t, err) {
						require.NoError(t, checkTable(g, 0))
					}
					cmd := CommandInput{
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

func (this *Helper) getOne(parts ...string) (doc resource.ObjectDocument, err error) {
	//"rooms", "lab", "contents"
	url := this.makeUrl(parts...)
	if resp, e := http.Get(url); e != nil {
		err = e
	} else {
		err = decodeBody(resp, &doc)
	}
	return
}

func (this *Helper) getMany(parts ...string) (doc resource.MultiDocument, err error) {
	//"rooms", "lab", "contents"
	url := this.makeUrl(parts...)
	if resp, e := http.Get(url); e != nil {
		err = e
	} else {
		err = decodeBody(resp, &doc)
	}
	return
}

func (this *Helper) post(input string) (doc resource.ObjectDocument, err error) {
	in := CommandInput{Input: input}
	return this.postCmd(in)
}

func (this *Helper) postCmd(in CommandInput) (doc resource.ObjectDocument, err error) {
	if b, e := json.Marshal(in); e != nil {
		err = e
	} else {
		postUrl := this.makeUrl()
		if resp, e := http.Post(postUrl, "application/json", bytes.NewReader(b)); e != nil {
			err = e
		} else if e := decodeBody(resp, &doc); e != nil {
			err = e
		} else {
			this.id = doc.Data.Id
		}
	}
	return doc, err
}

func (this *Helper) makeUrl(parts ...string) string {
	parts = append([]string{this.ts.URL, "game", this.id}, parts...)
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
