package commands

func endStory(output *CommandOutput) {
	output.NewCommand("finish", Dict{})
}
