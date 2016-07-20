package support

import (
	"github.com/ionous/sashimi/net/mem"
	"io/ioutil"
	"os"
	"os/user"
	"path"
)

const ext = ".sas"

// FileSaver implements MemSaver
// maybe add a file location?
type FileSaver struct {
	spath string // ex. "alice"
}

//saves:= "alice"
func NewFileSaver(spath string) *FileSaver {
	return &FileSaver{spath}
}

func (s *FileSaver) LoadBlob(slot string) (ret mem.SaveGameBlob, err error) {
	if where, e := s.savePath(); e != nil {
		err = e
	} else {
		name := path.Join(where, slot+ext)
		if f, e := os.Open(name); e != nil {
			err = e
		} else {
			defer f.Close()
			if i, e := f.Stat(); e != nil {
				err = e
			} else {
				b := make([]byte, i.Size())
				if _, e := f.Read(b); e != nil {
					err = e
				} else {
					ret = b
				}
			}
		}
	}
	return
}

// SaveBlob(slot string, blob SaveGameBlob) error
func (s *FileSaver) SaveBlob(_ string, b mem.SaveGameBlob) (slot string, err error) {
	if where, e := s.savePath(); e != nil {
		err = e
	} else if e := os.MkdirAll(where, 0777); e != nil {
		err = e
	} else {
		// the attitude of the go community can be very frustrating:
		// https://groups.google.com/forum/#!topic/golang-nuts/PHgye3Hm2_0
		if f, e := ioutil.TempFile(where, ""); e != nil {
			err = e
		} else if _, e := f.Write(b); e != nil {
			err = e
			f.Close()
		} else {
			src := f.Name()
			dst := src + ext
			if e := os.Rename(src, dst); e != nil {
				err = e
			} else {
				// name without extension
				slot = path.Base(src)
			}
		}
	}
	return
}

func (s *FileSaver) savePath() (where string, err error) {
	if usr, e := user.Current(); e != nil {
		err = e
	} else {
		where = path.Join(usr.HomeDir, s.spath)
	}
	return
}
