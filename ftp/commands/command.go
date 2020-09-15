package commands

type CommandType int

const (
	USER    CommandType = 0
	PASS    CommandType = 1
	SYST    CommandType = 2
	PWD     CommandType = 3
	CWD     CommandType = 4
	LIST    CommandType = 5
	TYPE    CommandType = 6
	EPRT    CommandType = 7
	LPRT    CommandType = 8
	UNKNOWN             = 500
)

type Command struct {
	Args    string
	CmdType CommandType
}

func ParseCommand(data []byte) Command {
	return GenParser(data).Eval()
}













