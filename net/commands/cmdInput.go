package commands

type CommandInput struct {
	Input   string `json:"in"`
	Action  string `json:"act"` // FIX? these should all probably be jsonapi refs
	Target  string `json:"tgt"`
	Context string `json:"ctx"`
}

func (this *CommandInput) Nouns() []string {
	nouns := make([]string, 0, 3)
	if n := this.Target; n != "" {
		nouns = append(nouns, n)
	} else if n := this.Context; n != "" {
		nouns = append(nouns, n)
	}
	return nouns
}
