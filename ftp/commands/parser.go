package commands

type Parser interface {
	Eval() Command
	next() Parser
}

func CommandParser(data []byte) Parser {
	return userParser{data: data}
}
