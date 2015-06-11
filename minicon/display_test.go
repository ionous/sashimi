package minicon

import (
	"strings"
	"testing"
)

//
func TestReflow(t *testing.T) {
	hl := Hardlines{}
	for _, w := range strings.Fields("Mary had a little lamb, which was stolen by aliens, and rescued by gundam.") {
		hl.AppendWord(w)
	}
	if len(hl.currentLine) != 14 {
		t.Fatal("expected 14 words", len(hl.currentLine))
	}
	lines := hl.Reflow(70)
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines %v", lines)
	}
}

//
func TestHistory(t *testing.T) {
	h := NewHistory(5)
	h.Add("old", nil)
	h.Add("mid", nil)
	h.Add("new", nil)
	rec := h.Mark()
	if was, ok := h.Back(); was != "new" || !ok {
		t.Fatal(was, "!= new")
	}
	if was, ok := h.Back(); was != "mid" || !ok {
		t.Fatal(was, "!= mid")
	}
	if was, ok := h.Back(); was != "old" || !ok {
		t.Fatal(was, "!= old")
	}
	if was, ok := h.Back(); was != "" || ok {
		t.Fatal(was, "!= ``; capped back")
	}
	if was, ok := h.Forward(); was != "mid" || !ok {
		t.Fatal(was, "!= mid; forward")
	}
	h.Restore(rec)
	if was, ok := h.Forward(); was != "" || ok {
		t.Fatalf("was:`%s` != ``, empty:%v, ok:%v", was, was != "", ok)
	}
	if was, ok := h.Back(); was != "new" || !ok {
		t.Fatal(was, "!= new; we didnt advance fwd, so we should now be back to the start of the test")
	}
}
