package commands

import (
	"strings"
)

type typeParser struct {
	data []byte
}

func (t typeParser) Eval() Command {
	command := string(t.data)

	if strings.HasPrefix(command, "TYPE") {
		return typeCommand()
	}

	return t.next().Eval()
}

func (t typeParser) next() Parser {
	return storParser{data: t.data}
}

func typeCommand() Command {
	return Command{CmdType: TYPE}
}
