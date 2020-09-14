package commands

import (
	"fmt"
	"net"
	"runtime"
	"strings"
)

type ConnectionManager struct {
	baseFs *FS
	acc    *accountManager
}

func NewConnectionManager(fs *FS) *ConnectionManager {
	return &ConnectionManager{baseFs: fs, acc: newAccountManager()}
}

func (connManager ConnectionManager) Handle(conn net.Conn) {
	defer conn.Close()

	conn.Write([]byte("220 Welcome to Nico FTP Server\n"))

	cmds := make(chan []byte)

	var buf [512]byte
	for {
		// read upto 512 bytes
		n, err := conn.Read(buf[0:])

		if err != nil {
			fmt.Println(err)
		}

		cmd := string(buf[0:n])

		fmt.Println("Received CMD " + cmd)

		go connManager.processCommand(buf[0:n], cmds)

		response := <-cmds

		fmt.Println("Sending: " + string(response))

		conn.Write(response)
	}
}

func (connManager ConnectionManager) processCommand(cmdData []byte, output chan []byte) {
	if !connManager.userCommand(cmdData, output) &&
		!connManager.passCommand(cmdData, output) &&
		!connManager.pwdCommand(cmdData, output) &&
		!connManager.systCommand(cmdData, output) &&
		!connManager.portCommand(cmdData, output) {

		output <- cmdData
	}

}

func (connManager ConnectionManager) userCommand(cmdData []byte, output chan []byte) bool {
	if cmd := string(cmdData); strings.HasPrefix(cmd, "USER") {
		user := strings.TrimSpace(cmd[5:])

		connManager.acc.withUser(user)

		// override the base virtual space with user specific virtual space
		connManager.baseFs = connManager.baseFs.ForUser(user)
		logMsg(user)

		response := "331 Need pass\n"
		logMsg(response)
		sendStr(response, output)

		return true
	}

	return false
}

func (connManager ConnectionManager) passCommand(cmdData []byte, output chan []byte) bool {
	if cmd := string(cmdData); strings.HasPrefix(cmd, "PASS") {
		pass := strings.TrimSpace(cmd[4:])
		if connManager.acc.validatePassword(pass) {
			response := "230 User logged in, proceed.\n"

			logMsg(response)
			sendStr(response, output)
		} else {
			response := "530 Incorrect Pass.\n"

			logMsg(response)
			sendStr(response, output)
		}

		return true
	}

	return false
}

func (connManager ConnectionManager) systCommand(cmdData []byte, output chan []byte) bool {
	if cmd := string(cmdData); strings.HasPrefix(cmd, "SYST") {
		sysName := runtime.GOOS

		logMsg(fmt.Sprintf("SYSTEM NAME: %v", sysName))

		response := fmt.Sprintf("215 TYPE: %v", sysName)
		logMsg(response)

		sendStr(response, output)

		return true
	}

	return false
}

func (connManager ConnectionManager) pwdCommand(cmdData []byte, output chan []byte) bool {
	if cmd := string(cmdData); strings.HasPrefix(cmd, "PWD") {
		response := fmt.Sprintf("257 %v\r\n", "/Users/anicolaspp")

		logMsg(response)

		sendStr(response, output)

		return true
	}

	return false
}

func (connManager ConnectionManager) portCommand(cmdData []byte, output chan []byte) bool {
	if cmd := string(cmdData); strings.HasPrefix(cmd, "PORT") {
		logMsg(cmd)

		sendStr(cmd, output)

		return true
	}

	return false
}

func logMsg(value interface{}) {
	fmt.Println(value)
}

func sendStr(msg string, to chan<- []byte) {
	to <- []byte(msg)
}
