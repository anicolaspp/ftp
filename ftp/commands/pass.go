package commands

import (
	"strings"
)

type passParser struct {
	data []byte
}

func (p passParser) Eval() Command {
	command := string(p.data)

	if strings.HasPrefix(command, "PASS") {
		return passCommand(strings.TrimSpace(command[5:]))
	}

	return p.next().Eval()
}

func (p passParser) next() Parser {
	return sysParser{data: p.data}
}

func passCommand(args string) Command {
	return Command{Args: args, CmdType: PASS}
}