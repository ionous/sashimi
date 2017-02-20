package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ionous/mars/export"
	"github.com/ionous/mars/facts"
	"github.com/ionous/mars/std"
	"github.com/ionous/mars/tools/uniform"
	"github.com/ionous/sashimi/_examples/stories"
	"github.com/ionous/sashimi/standard/output"
	"log"
	"os"
)

func Marshall(src interface{}) (ret string, err error) {
	b := new(bytes.Buffer)
	enc := json.NewEncoder(b)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", " ")
	if e := enc.Encode(src); e != nil {
		err = e
	} else {
		ret = b.String()
	}
	return
}

func write(dst, m string) {
	log.Println("writing to", dst)
	if f, e := os.Create(dst); e != nil {
		log.Println("error", e)
	} else {
		defer f.Close()
		fmt.Fprintln(f, m)
	}
}

// Boilerplate to run the story.
// go run play.go -export -story empty -file /Users/ionous/Dev/makisu/bin/empty.js
// go run play.go -export -story sushi -file /Users/ionous/Dev/makisu/bin/sushi.js
func main() {
	storyName := flag.String("story", "", "select the story to play.")
	exportFlag := flag.Bool("export", false, "true to export story.")
	fileName := flag.String("file", "", "export destination.")
	options := output.ParseCommandLine()

	if s, ok := stories.Select(*storyName); !ok {
		flag.PrintDefaults()
		fmt.Println("Please select one of the following stories:")
		for _, nick := range stories.List() {
			fmt.Println(" ", nick)
		}
	} else if !*exportFlag {
		output.RunGame(*s, options)
	} else {
		ctx := uniform.NewContext()
		if sections, e := uniform.NewLibraries(ctx, facts.Facts(), export.Export(), std.Std()); e != nil {
			fmt.Println("libraries error", e)
		} else if chapter, e := uniform.NewChapter("Chapter One", s.Declarations()); e != nil {
			fmt.Println("chapter error", e)
		} else {
			story := export.Story{*storyName, append(sections, chapter)}
			enc := uniform.NewUniformEncoder(ctx.Types)
			if data, e := enc.Compute(story); e != nil {
				fmt.Println("compute error", e)
			} else if m, e := Marshall(data); e != nil {
				fmt.Println("marshal error", e)
			} else {
				if *fileName == "" {
					fmt.Println(m)
				} else {
					write(*fileName+"on", m)
					write(*fileName, "module.exports="+m+";")
				}
			}
		}
	}
}
