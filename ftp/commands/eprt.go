package commands

import (
	"strings"
)

type eprtParser struct {
	data []byte
}
func (e eprtParser) Eval() Command {
	command := string(e.data)

	if strings.HasPrefix(command, "EPRT") {
		return eprtCommand(strings.TrimSpace(command[5:]))
	}

	return e.next().Eval()
}

func (e eprtParser) next() Parser {
	return  lprtParser{data: e.data}
}

func eprtCommand(args string) Command {
	return Command{Args: args, CmdType: EPRT}
}