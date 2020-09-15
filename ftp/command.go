package ftp

import "strings"

func ParseCommand(data []byte) Command {
	command := string(data)

	if strings.HasPrefix(command, "USER") {
		return userCommand(strings.TrimSpace(command[5:]))
	} else if strings.HasPrefix(command, "PASS") {
		return passCommand(strings.TrimSpace(command[5:]))
	} else if strings.HasPrefix(command, "SYST") {
		return systCommand(strings.TrimSpace(command[5:]))
	} else if strings.HasPrefix(command, "PWD") {
		return pwdCommand()
	} else if strings.HasPrefix(command, "CWD") {
		return cwdCommand(strings.TrimSpace(command[4:]))
	} else if strings.HasPrefix(command, "LIST") {
		return lsCommand()
	} else if strings.HasPrefix(command, "EPRT") {
		return eprtCommand(strings.TrimSpace(command[5:]))
	} else if strings.HasPrefix(command, "LPRT") {
		return lprtCommand(strings.TrimSpace(command[5:]))
	} else if strings.HasPrefix(command, "TYPE") {
		return typeCommand()
	}

	return Command{cmdType: UNKNOWN}
}

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
	cmdType CommandType
}

func userCommand(args string) Command {
	return Command{Args: args, cmdType: USER}
}

func passCommand(args string) Command {
	return Command{Args: args, cmdType: PASS}
}

func systCommand(args string) Command {
	return Command{Args: args, cmdType: SYST}
}

func pwdCommand() Command {
	return Command{cmdType: PWD}
}

func cwdCommand(args string) Command {
	return Command{Args: args, cmdType: CWD}
}

func lsCommand() Command {
	return Command{cmdType: LIST}
}

func lprtCommand(args string) Command {
	return Command{Args: args, cmdType: LPRT}
}

func eprtCommand(args string) Command {
	return Command{Args: args, cmdType: EPRT}
}

func typeCommand() Command {
	return Command{cmdType: TYPE}
}
