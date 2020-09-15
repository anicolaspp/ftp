package commands

import (
	"strings"
)

type cwdParser struct {
	data []byte
}

func (c cwdParser) Eval() Command {
	command := string(c.data)

	if strings.HasPrefix(command, "CWD") {
		return cwdCommand(strings.TrimSpace(command[4:]))
	}

	return c.next().Eval()
}

func (c cwdParser) next() Parser {
	return lsParser{data: c.data}
}

func cwdCommand(args string) Command {
	return Command{Args: args, CmdType: CWD}
}
