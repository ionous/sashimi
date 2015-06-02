package script

import (
	S "github.com/ionous/sashimi/source"
)

type ListOfItems struct {
	origin Origin
	list   string
	items  []string
}

func (this ListOfItems) And(name string) ListOfItems {
	this.items = append(this.items, name)
	return this
}

// FIX: move these into a standard rules extension package?
func In(room string) ListOfItems {
	return ListOfItems{NewOrigin(2), "whereabouts", []string{room}}

}

func Supports(prop string) ListOfItems {
	return ListOfItems{NewOrigin(2), "contents", []string{prop}}
}

func Contains(prop string) ListOfItems {
	return ListOfItems{NewOrigin(2), "contents", []string{prop}}
}

func Possesses(prop string) ListOfItems {
	return ListOfItems{NewOrigin(2), "inventory", []string{prop}}
}

func Wears(prop string) ListOfItems {
	return ListOfItems{NewOrigin(2), "clothing", []string{prop}}
}

func (this ListOfItems) MakeStatement(b SubjectBlock) (err error) {
	list, code := this.list, this.origin.Code()
	for _, item := range this.items {
		fields := S.KeyValueFields{b.subject, list, item}
		if e := b.NewKeyValue(fields, code); e != nil {
			err = e
			break
		}
	}
	return err
}
