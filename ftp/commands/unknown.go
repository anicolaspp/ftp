package commands

type unknownParser int

func (u unknownParser) Eval() Command {
	return Command{CmdType: UNKNOWN}
}

func (u unknownParser) next() Parser {
	return nil
}