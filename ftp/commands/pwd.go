package commands

import (
	"strings"
)

type pwdParser struct {
	data []byte
}

func (p pwdParser) Eval() Command {
	command := string(p.data)

	if strings.HasPrefix(command, "PWD") {
		return pwdCommand()
	}

	return p.next().Eval()
}

func (p pwdParser) next() Parser {
	return cwdParser{data: p.data}
}

func pwdCommand() Command {
	return Command{CmdType: PWD}
}

