package commands

import "strings"

type quitParser struct {
	data []byte
}


func (q quitParser) Eval() Command {
	command := string(q.data)

	if strings.HasPrefix(command, "QUIT") {
		return Command{CmdType: QUIT}
	}

	return q.next().Eval()
}

func (q quitParser) next() Parser {
	return new(unknownParser)
}
