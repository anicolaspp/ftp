package commands

import (
	"strings"
)

type userParser struct {
	data []byte
}

func (u userParser) Eval() Command {
	command := string(u.data)

	if strings.HasPrefix(command, "USER") {
		return userCommand(strings.TrimSpace(command[5:]))
	}

	return u.next().Eval()
}

func (u userParser) next() Parser {
	return passParser{data: u.data}
}

func userCommand(args string) Command {
	return Command{Args: args, CmdType: USER}
}