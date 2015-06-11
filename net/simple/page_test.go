package simple

import (
	"os"
	"testing"
)

func TestWebTemplate(t *testing.T) {
	lines := []string{"here's a line"}
	if e := page.ExecuteTemplate(os.Stdout, "simple.html", lines); e != nil {
		t.Fatal(e)
	}
}
