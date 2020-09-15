package commands

import (
	"strings"
)

type lprtParser struct {
	data []byte
}

func (l lprtParser) Eval() Command {
	command := string(l.data)

	if strings.HasPrefix(command, "LPRT") {
		return lprtCommand(strings.TrimSpace(command[5:]))
	}

	return l.next().Eval()
}

func (l lprtParser) next() Parser {
	return typeParser{data: l.data}
}

func lprtCommand(args string) Command {
	return Command{Args: args, CmdType: LPRT}
}

