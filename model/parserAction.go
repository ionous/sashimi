package model

// map actions to the parser
type ParserAction struct {
	action   *ActionInfo
	commands []string
}

func NewParserAction(action *ActionInfo, commands []string) ParserAction {
	return ParserAction{action, commands}
}

func (this *ParserAction) Action() *ActionInfo {
	return this.action
}

func (this *ParserAction) Commands() []string {
	return this.commands
}
