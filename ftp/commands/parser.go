package commands

type Parser interface {
	Eval() Command
	next() Parser
}

func GenParser(data []byte) Parser {
	return userParser{data: data}
}
