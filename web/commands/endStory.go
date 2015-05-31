package commands

func endStory(output *CommandOutput) {
	output.Add("finish", Dict{})
}
