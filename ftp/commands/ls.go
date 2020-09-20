package commands

import (
	"strings"
)

type lsParser struct {
	data []byte
}

func (l lsParser) Eval() Command {
	command := string(l.data)

	if strings.HasPrefix(command, "LIST") {
		return lsCommand()
	}

	return l.next().Eval()
}

func (l lsParser) next() Parser {
	return eprtParser{data: l.data}
}


func lsCommand() Command {
	return Command{CmdType: LIST}
}