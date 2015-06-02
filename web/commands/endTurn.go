package commands

func endTurn(output *CommandOutput) {
	output.NewCommand("turn", Dict{})
}
