package commands

import "strings"

type retrParser struct {
	data []byte
}

func (r retrParser) Eval() Command {
	if command := string(r.data); strings.HasPrefix(command, "RETR") {
		return Command{Args: strings.TrimSpace(command[5:]), CmdType: RETR}
	}

	return r.next().Eval()
}

func (r retrParser) next() Parser {
	return quitParser{data: r.data}
}
