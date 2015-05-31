package commands

func endTurn(output *CommandOutput) {
	output.Add("turn", Dict{})
}
