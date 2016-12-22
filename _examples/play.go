package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ionous/mars/encode"
	"github.com/ionous/sashimi/_examples/stories"
	"github.com/ionous/sashimi/standard/output"
	"os"
)

// Boilerplate to run the story.
func main() {
	story := flag.String("story", "", "select the story to play.")
	export := flag.Bool("export", false, "true to export story.")
	dst := flag.String("file", "", "export destination.")
	options := output.ParseCommandLine()

	if s, ok := stories.Select(*story); !ok {
		flag.PrintDefaults()
		fmt.Println("Please select one of the following stories:")
		for _, nick := range stories.List() {
			fmt.Println(" ", nick)
		}
	} else if !*export {
		output.RunGame(s, options)
	} else if cmd, e := encode.Compute(*s); e != nil {
		fmt.Println("error", e)
	} else if b, e := json.MarshalIndent(cmd, "", " "); e != nil {
		fmt.Println("error", e)
	} else {
		w := os.Stdout
		if *dst != "" {
			fmt.Println("writing to", *dst)
			if f, e := os.Create(*dst); e != nil {
				fmt.Println("error", e)
				return
			} else {
				w = f
				defer f.Close()
			}
		}
		fmt.Fprintln(w, string(b))
		// fmt.Println(w, string(b))
	}
}
