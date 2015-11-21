package event

import (
	"fmt"
	"github.com/ionous/sashimi/util/ident"
)

var _ = fmt.Sprintf

type Message struct {
	Id          ident.Id
	Data        interface{}
	CaptureOnly bool
	CantCancel  bool
}

func (msg *Message) String() string {
	return msg.Id.String()
}

// Send returns true if the default action is desired.
func (msg *Message) Send(path PathList) (bool, error) {
	okay := true
	if path.Len() > 0 {
		// data:
		target := path.Cast(path.Back())
		proc := &Proc{msg: msg, path: path, target: target}

		// capture, all the way down to, but not including, the target
		if !proc.stopMore {
			proc.phase = CapturingPhase
			for it := path.Front(); it != path.Back(); it = it.Next() {
				loc := path.Cast(it)
				//fmt.Println("event.Messge", proc.phase, loc)
				if e := proc.sendToTarget(loc); e != nil {
					return false, e
				}
				if proc.stopNow {
					break
				}
			}
		}

		// target
		if !proc.stopMore {
			proc.phase = TargetPhase
			loc := proc.target
			//fmt.Println("event.Messge", proc.phase, loc)
			if e := proc.sendToTarget(loc); e != nil {
				return false, e
			}
		}

		// bubble, all the way to, and including, the root
		if !proc.stopMore {
			proc.phase = BubblingPhase
			for it := path.Back().Prev(); it != nil; it = it.Prev() {
				loc := path.Cast(it)
				//fmt.Println("event.Messge", proc.phase, loc)
				if e := proc.sendToTarget(loc); e != nil {
					return false, e
				}
				if proc.stopNow {
					break
				}
			}
		}
		okay = !proc.cancelled
	}
	return okay, nil
}
