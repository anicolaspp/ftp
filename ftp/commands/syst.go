package commands

import (
	"strings"
)

type sysParser struct {
	data []byte
}

func (s sysParser) Eval() Command {
	command := string(s.data)

	if strings.HasPrefix(command, "SYST") {
		return systCommand(strings.TrimSpace(command[5:]))
	}

	return s.next().Eval()
}

func (s sysParser) next() Parser {
	return pwdParser{data: s.data}
}

func systCommand(args string) Command {
	return Command{Args: args, CmdType: SYST}
}