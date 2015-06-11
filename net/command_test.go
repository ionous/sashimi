package net

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ionous/sashimi/_examples/stories"
	"github.com/ionous/sashimi/net/commands"
	"github.com/ionous/sashimi/net/resource"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	bodyType = "application/json"
)

//
func TestCommandGame(t *testing.T) {
	stories.Select("lab")

	ts := httptest.NewServer(commands.NewServer())
	defer ts.Close()

	g := &Helper{ts, "new"}
	if d, err := g.post(""); assert.NoError(t, err) {
		if assert.Len(t, d.Included, 2, "the player and the room") {

			if changes, ok := d.Data.Attributes["events"]; assert.True(t, ok, "frame has event stream") {
				changes := changes.([]interface{})
				assert.True(t, len(changes) > 1)

				assert.EqualValues(t, "game", d.Data.Class)
				// check the room
				if contents, err := g.get("rooms", "lab", "contents"); assert.NoError(t, err) {
					assert.Len(t, contents.Data, 2, "the lab should have two objects")
					assert.Len(t, contents.Included, 1, "the player should be previously known, the table newly known.")
				}
				checkTable(t, g, 1)

				if _, err := g.post("open the glass jar"); assert.NoError(t, err) {
					checkTable(t, g, 1)
				}

				// take the beaker
				if _, err := g.post("take the glass jar"); assert.NoError(t, err) {
					checkTable(t, g, 0)
				}
			}
		}
	}
}

func checkTable(t *testing.T, g *Helper, cnt int) {
	if contents, err := g.get("supporters", "table", "contents"); assert.NoError(t, err) {
		t.Log(pretty(contents))
		assert.Len(t, contents.Data, cnt, "the table should have %d objects", cnt)
	}
}

type Helper struct {
	ts *httptest.Server
	id string
}

func (this *Helper) get(parts ...string) (doc resource.MultiDocument, err error) {
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
	in := commands.CommandInput{Input: input}
	if b, e := json.Marshal(in); e != nil {
		err = e
	} else {
		postUrl := this.makeUrl()
		if resp, e := http.Post(postUrl, bodyType, bytes.NewReader(b)); e != nil {
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
