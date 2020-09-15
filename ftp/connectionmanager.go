package ftp

import (
	"encoding/binary"
	"fmt"
	"github.com/anicolaspp/ftp/ftp/commands"
	"net"
	"runtime"
	"strconv"
	"strings"
)

type ConnectionManager struct {
	fs  *FS
	acc *accountManager

	ctlConnection net.Conn

	dataConnection net.Conn
}

func NewConnectionManager(baseDir string) *ConnectionManager {
	return &ConnectionManager{acc: newAccountManager(), fs: NewFS(baseDir)}
}

func (connManager *ConnectionManager) Handle(conn net.Conn) {
	connManager.ctlConnection = conn

	defer connManager.ctlConnection.Close()

	connManager.ctlConnection.Write([]byte("220 Welcome to Nico FTP Server\n"))

	// start connection control loop
	var buf [512]byte

	for {
		// read upto 512 bytes
		n, err := conn.Read(buf[0:])

		if err != nil {
			fmt.Println(err)
			return
		}

		cmd := string(buf[0:n])

		logMsg(fmt.Sprintf("[CLIENT CMD]: %v\n", cmd))

		connManager.processCommand(buf[0:n])
	}
}

func (connManager *ConnectionManager) processCommand(cmdData []byte) {
	cmd := commands.ParseCommand(cmdData)

	fmt.Println(cmd)

	if !connManager.user(cmd) &&
		!connManager.pass(cmd) &&
		!connManager.pwd(cmd) &&
		!connManager.syst(cmd) &&
		!connManager.lprt(cmd) &&
		!connManager.typ(cmd) &&
		!connManager.eprt(cmd) &&
		!connManager.list(cmd) &&
		!connManager.cwd(cmd) {

		connManager.echo(cmdData)
	}

}

func (connManager *ConnectionManager) echo(cmdData []byte) {
	connManager.sendStr(fmt.Sprintf("%v\n", string(cmdData)))
}

func (connManager *ConnectionManager) user(cmd commands.Command) bool {

	if cmd.CmdType == commands.USER {
		user := cmd.Args

		connManager.acc.withUser(user)

		// override the base virtual space with user specific virtual space
		connManager.fs = connManager.fs.ForUser(user)
		logMsg(user)

		response := "331 Need pass\n"
		connManager.sendStr(response)

		return true
	}

	return false
}

func (connManager *ConnectionManager) pass(cmd commands.Command) bool {
	if cmd.CmdType == commands.PASS {
		pass := cmd.Args

		if connManager.acc.validatePassword(pass) {
			response := "230 User logged in, proceed.\n"

			connManager.sendStr(response)
		} else {
			response := "530 Incorrect Pass.\n"

			connManager.sendStr(response)
		}

		return true
	}

	return false
}

func (connManager *ConnectionManager) syst(cmd commands.Command) bool {
	if cmd.CmdType == commands.SYST {
		sysName := runtime.GOOS

		response := fmt.Sprintf("215 TYPE: %v\n", sysName)

		connManager.sendStr(response)

		return true
	}

	return false
}

func (connManager *ConnectionManager) pwd(cmd commands.Command) bool {
	if cmd.CmdType == commands.PWD {
		response := fmt.Sprintf("257 %v\n", connManager.fs.Pwd())

		connManager.sendStr(response)

		return true
	}

	return false
}

func (connManager *ConnectionManager) cwd(cmd commands.Command) bool {
	if cmd.CmdType == commands.CWD {
		path := cmd.Args

		currentPath, err := connManager.fs.Cwd(path)

		if err != nil {
			errorResponse := fmt.Sprintf("550 %v\n", err)

			connManager.sendStr(errorResponse)
		} else {
			response := fmt.Sprintf("250 OK. Change path to %v\n", currentPath)

			connManager.sendStr(response)
		}

		return true
	}

	return false
}

func (connManager *ConnectionManager) list(cmd commands.Command) bool {
	if cmd.CmdType == commands.LIST {

		connManager.sendStr("150 Listing Directory Content\n")

		content := connManager.fs.Ls()

		strResponse := strings.Join(content, "\n") + "\n"

		connManager.dataConnection.Write([]byte(strResponse))
		connManager.dataConnection.Close()

		connManager.sendStr("226 Directory send OK\r\n")

		return true
	}

	return false
}

func (connManager *ConnectionManager) eprt(cmd commands.Command) bool {
	if cmd.CmdType == commands.EPRT {
		args := strings.Split(strings.Trim(cmd.Args, "\r\n"), "|")

		port, _ := strconv.Atoi(args[3])

		if connManager.openDataConnection(nil, int64(port)) {
			response := "200 Get Port\n"
			connManager.sendStr(response)
		}

		return true
	}

	return false
}

func (connManager *ConnectionManager) lprt(cmd commands.Command) bool {
	if cmd.CmdType == commands.LPRT {

		args := strings.Split(strings.Trim(cmd.Args, "\r\n"), ",")

		addSize, _ := strconv.Atoi(args[1])

		add := make([]byte, addSize)

		for i := 0; i < addSize; i++ {
			v, _ := strconv.Atoi(args[2+i])

			add[i] = byte(v)
		}

		fmt.Println(addSize)
		fmt.Println(add)

		portSize, _ := strconv.Atoi(args[2+addSize])

		port := make([]byte, portSize)

		for i := 0; i < portSize; i++ {
			v, _ := strconv.Atoi(args[3+addSize+i])

			port[i] = byte(v)
		}

		if portSize < 8 {
			for i := 0; i < 8-portSize; i++ {
				port = append([]byte{0}, port...)
			}
		}

		portLong := binary.BigEndian.Uint64(port)

		var ip net.IP = add

		fmt.Println(ip)
		fmt.Println(portLong)

		if connected := connManager.openDataConnection(ip, int64(portLong)); connected {
			response := "200 Get Port\n"
			connManager.sendStr(response)
		} else {

		}

		return true
	}

	return false
}

func (connManager *ConnectionManager) typ(cmd commands.Command) bool {
	if cmd.CmdType == commands.TYPE {
		response := "200\n"

		connManager.sendStr(response)

		return true
	}

	return false
}

func (connManager *ConnectionManager) sendStr(msg string) {

	logMsg(fmt.Sprintf("[SERVER]: %v\n", msg))

	connManager.ctlConnection.Write([]byte(msg))
}

func (connManager *ConnectionManager) openDataConnection(ip net.IP, port int64) bool {

	address := fmt.Sprintf("%v:%v", "localhost", port)

	logMsg(fmt.Sprintf("[SERVER]: Connecting to %v\n", address))

	dataConn, error := net.Dial("tcp", address)

	if error != nil {
		fmt.Printf("Error opening data connection: %v", error)
		return false
	}

	logMsg(fmt.Sprintf("Data connection opened at %v\n", address))

	connManager.dataConnection = dataConn

	return true
}

func logMsg(value interface{}) {
	fmt.Print(value)
}
