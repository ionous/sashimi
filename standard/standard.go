package standard

import (
	"fmt"
	G "github.com/ionous/sashimi/game"
	M "github.com/ionous/sashimi/model"
	R "github.com/ionous/sashimi/runtime"
	"log"
)

type StandardCore struct {
	*R.Game
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
			storyObject := game.FindFirstOf(model.Classes.FindClass("stories"))
			if storyObject == nil {
				err = fmt.Errorf("couldn't find story")
			} else {
				if statusObject, ok := game.FindObject("status bar"); !ok {
					err = fmt.Errorf("couldn't find status bar")
				} else {
					story := R.NewObjectAdapter(game, storyObject)
					status := R.NewObjectAdapter(game, statusObject)
					//
					core := StandardCore{game, parser, output, story, status}
					core.SetLeft(story.Name())
					core.SetRight(fmt.Sprint(story.Name(), "by ", story.Text("author")))
					ret = StandardStart{core}
				}
			}
		}
	}
	return ret, err
}

//
// sends starting to play, and returns a new game.
//
func (this *StandardStart) Start() (ret StandardGame, err error) {
	// FIX: shouldnt the interface be Go("commence")?
	if e := this.SendEvent("starting to play", this.story.String()); e != nil {
		err = e
	} else {
		// process all existing messages in the queue first
		if e := this.ProcessEvents(); e != nil {
			err = e
		}
	}
	return StandardGame{this.StandardCore, false, false}, err
}

func (this *StandardGame) IsQuit() bool {
	return this.quit
}

func (this *StandardGame) IsFinished() bool {
	return this.quit || this.completed
}

//
// return false if the game has finished
//
func (this *StandardGame) Input(s string) bool {
	if !this.IsFinished() {
		game, out, parser := this.Game, this.output, this.parser
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
	return !this.IsFinished()
}
