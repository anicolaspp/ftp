package ftp

import (
	"encoding/binary"
	"fmt"
	"net"
	"runtime"
	"strconv"
	"strings"
)

type ConnectionManager struct {
	baseFs *FS
	acc    *accountManager

	ctlConnection net.Conn

	dataConnection net.Conn
}

func NewConnectionManager(fs *FS) *ConnectionManager {
	return &ConnectionManager{baseFs: fs, acc: newAccountManager()}
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
	if !connManager.user(cmdData) &&
		!connManager.pass(cmdData) &&
		!connManager.pwd(cmdData) &&
		!connManager.syst(cmdData) &&
		!connManager.lprt(cmdData) &&
		!connManager.typ(cmdData) &&
		!connManager.eprt(cmdData) &&
		!connManager.list(cmdData) {

		connManager.echo(cmdData)
	}

}

func (connManager *ConnectionManager) echo(cmdData []byte) {
	connManager.sendStr(fmt.Sprintf("%v\n", string(cmdData)))
}

func (connManager *ConnectionManager) user(cmdData []byte) bool {
	if cmd := string(cmdData); strings.HasPrefix(cmd, "USER") {
		user := strings.TrimSpace(cmd[5:])

		connManager.acc.withUser(user)

		// override the base virtual space with user specific virtual space
		connManager.baseFs = connManager.baseFs.ForUser(user)
		logMsg(user)

		response := "331 Need pass\n"
		connManager.sendStr(response)

		return true
	}

	return false
}

func (connManager *ConnectionManager) pass(cmdData []byte) bool {
	if cmd := string(cmdData); strings.HasPrefix(cmd, "PASS") {
		pass := strings.TrimSpace(cmd[4:])
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

func (connManager *ConnectionManager) syst(cmdData []byte) bool {
	if cmd := string(cmdData); strings.HasPrefix(cmd, "SYST") {
		sysName := runtime.GOOS

		response := fmt.Sprintf("215 TYPE: %v\n", sysName)

		connManager.sendStr(response)

		return true
	}

	return false
}

func (connManager *ConnectionManager) pwd(cmdData []byte) bool {
	if cmd := string(cmdData); strings.HasPrefix(cmd, "PWD") {
		response := fmt.Sprintf("257 %v\n", "\"/Users/anicolaspp\"")

		connManager.sendStr(response)

		return true
	}

	return false
}

func (connManager *ConnectionManager) list(cmdData []byte) bool {
	if cmd := string(cmdData); strings.HasPrefix(cmd, "LIST") {

		connManager.sendStr("150 Listing Directory Content\n")

		connManager.dataConnection.Write([]byte("my_dir\n"))
		connManager.dataConnection.Close()

		connManager.sendStr("226 Directory send OK\r\n")

		return true
	}

	return false
}

func (connManager *ConnectionManager) eprt(cmdData []byte) bool {
	if cmd := string(cmdData); strings.HasPrefix(cmd, "EPRT") {
		args := strings.Split(strings.Trim(cmd[5:], "\r\n"), "|")

		port, _ := strconv.Atoi(args[3])

		if connManager.openDataConnection(nil, int64(port)) {
			response := "200 Get Port\n"
			connManager.sendStr(response)
		}

		return true
	}

	return false
}

func (connManager *ConnectionManager) lprt(cmdData []byte) bool {
	if cmd := string(cmdData); strings.HasPrefix(cmd, "LPRT") {

		args := strings.Split(strings.Trim(cmd[5:], "\r\n"), ",")

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

func (connManager *ConnectionManager) typ(cmdData []byte) bool {
	if cmd := string(cmdData); strings.HasPrefix(cmd, "TYPE") {
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
