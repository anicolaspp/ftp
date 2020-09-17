package commands

import "strings"

type storParser struct {
	data []byte
}

func (s storParser) Eval() Command {
	command := string(s.data)

	if strings.HasPrefix(command, "STOR") {
		return Command{Args: strings.TrimSpace(command[5:]), CmdType: STOR}
	}

	return s.next().Eval()
}

func (s storParser) next() Parser {
	return quitParser{data: s.data}
}
