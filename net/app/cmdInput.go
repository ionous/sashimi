package app

// CommandInput defines the expected post data format from a client.
// ie. this aint REST.
type CommandInput struct {
	// Input accepts a raw input string from the user: when specified, the other fields are ignored.
	Input string `json:"in"`
	// Action accepts all known player actions: the source of player commands is always the player.
	Action string `json:"act"`
	// Target: the first noun of the action
	Target string `json:"tgt"`
	// Context: the second noun of the action.
	Context string `json:"ctx"`
}

// Nouns returns an array of target and context strings, skipping empty strings.
func (cmd *CommandInput) Nouns() []string {
	nouns := make([]string, 0, 3)
	if n := cmd.Target; n != "" {
		nouns = append(nouns, n)
	}
	if n := cmd.Context; n != "" {
		nouns = append(nouns, n)
	}
	return nouns
}
