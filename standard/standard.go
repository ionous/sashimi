package standard

import (
	"fmt"
	G "github.com/ionous/sashimi/game"
	M "github.com/ionous/sashimi/model"
	//P "github.com/ionous/sashimi/parser"
	R "github.com/ionous/sashimi/runtime"
	"log"
)

type StandardCore struct {
	game          *R.Game
	parser        *R.ModelParser
	output        R.IOutput
	story, status R.ObjectAdapter
}
type StandardStart struct {
	StandardCore
}
type StandardGame struct {
	StandardCore
	quit, completed bool
}

func (this *StandardCore) Left() string {
	return this.status.Text("left")
}

func (this *StandardCore) Right() string {
	return this.status.Text("right")
}

func (this *StandardCore) SetLeft(status string) {
	this.status.SetText("left", status)
}

func (this *StandardCore) SetRight(status string) {
	this.status.SetText("right", status)
}

func NewStandardGame(model *M.Model, output R.IOutput) (ret StandardStart, err error) {
	if game, e := R.NewGame(model, output); e != nil {
		err = e
	} else {
		if parser, e := R.NewParser(game); e != nil {
			err = e
		} else {
			// FIX: move parser source into the model/parser
			game.PushParserSource(func(g G.Play) (ret G.IObject) {
				return g.The("player")
			})
			game.PushParentLookup(func(g G.Play, o G.IObject) (ret G.IObject) {
				if parent, where := DirectParent(o); where != "" {
					ret = parent
				}
				return ret
			})
			story := game.FindFirstOf(model.Classes.FindClass("stories"))
			if story == nil {
				err = fmt.Errorf("couldn't find story")
			} else {
				if obj, okay := game.FindObject("status bar"); !okay {
					err = fmt.Errorf("couldn't find status bar")
				} else {
					story := R.NewObjectAdapter(game, story)
					status := R.NewObjectAdapter(game, obj)
					ret = StandardStart{StandardCore{game, parser, output, story, status}}
				}
			}
		}
	}
	return ret, err
}

func (this *StandardStart) Start() (ret StandardGame, err error) {
	// FIX: shouldnt the interface be Go("commence")?
	if e := this.game.SendEvent("starting to play", this.story.String()); e != nil {
		err = e
	} else {
		// process all existing messages in the queue first
		if e := this.game.ProcessEvents(); e != nil {
			err = e
		}
	}
	return StandardGame{this.StandardCore, false, false}, err
}

func (this *StandardGame) Quit() bool {
	return this.quit
}

func (this *StandardGame) Finished() bool {
	return this.quit || this.completed
}

func (this *StandardGame) Input(s string) bool {
	if !this.Finished() {
		game, out, parser := this.game, this.output, this.parser
		in := parser.NormalizeInput(s)
		if in == "q" || in == "quit" {
			this.quit = true
		} else {
			if _, res, e := parser.Parse(in); e != nil {
				out.Println(e)
			} else if e := res.Run(); e != nil {
				out.Println(e)
			}
			game.SendEvent("ending the turn", this.story.String())
			if e := game.ProcessEvents(); e != nil {
				log.Println(e)
			} else {
				if this.story.Is("completed") {
					this.completed = true
				}
			}
		}
	}
	return !this.Finished()
}
