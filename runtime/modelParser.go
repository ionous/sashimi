package runtime

import (
	M "github.com/ionous/sashimi/model"
	P "github.com/ionous/sashimi/parser"
)

//
// adapt the parser to the model
//
type ModelParser struct {
	*P.Parser
}

func (this *ModelParser) NormalizeInput(s string) string {
	return P.NormalizeInput(s)
}

//
// create the parser, add all commands and patterns
//
func NewParser(game *Game,
) (ret *ModelParser, err error,
) {
	model := game.Model
	parser := P.NewParser()
	// pre-compile the parser statements ( ex. to catch errors. )
CreateActions:
	for _, p := range model.ParserActions {
		act := p.Action()
		cmd := ParserCommand{act, model, game}
		// we expect the game to supply the first noun of every action
		comprehension, e := parser.AddCommand(act.Action(), cmd, act.NumNouns()-1)
		if e != nil {
			err = e
			break CreateActions
		}

		for _, learn := range p.Commands() {
			if e := comprehension.LearnPattern(learn); e != nil {
				err = e
				break CreateActions
			}
		}
	}
	if err == nil {
		ret = &ModelParser{parser}
	}
	return ret, err
}

//
// implements P.ICommand
//
type ParserCommand struct {
	act   *M.ActionInfo
	model *M.Model
	game  *Game
}

func (this ParserCommand) NewMatcher() P.IMatch {
	return &NounFactory{this.act, this.model, 0}
}

// our matcher, the noun factory, yields nouns in the form of instance string ids
func (this ParserCommand) RunCommand(nouns ...string) (err error) {
	return this.game.RunAction(this.act, nouns)
}
