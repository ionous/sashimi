package play

import (
	E "github.com/ionous/sashimi/event"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/play/api"
	"github.com/ionous/sashimi/util/sbuf"
	"io"
	"strings"
)

type noFrame struct {
	log   io.Writer
	parts []string
}

func (d *noFrame) BeginEvent(_, _ meta.Instance, path E.PathList, msg *E.Message) api.IEndEvent {
	d.parts = append(d.parts, msg.String())
	fullName := strings.Join(d.parts, "/")
	d.log.Write([]byte(sbuf.New("sending", fullName, "to", path).Join(" ")))
	return d
}

func (d *noFrame) FlushFrame() {
}

func (d *noFrame) EndEvent() {
	d.parts = d.parts[0 : len(d.parts)-1]
}
