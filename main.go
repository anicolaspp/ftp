package main

import (
	"fmt"
	"github.com/anicolaspp/ftp/commands"
	"net"
	"os"
)

func main() {

	listen, err := net.Listen("tcp", ":21")

	checkError(err)

	defer listen.Close()

	for {
		conn, err := listen.Accept()
		checkError(err)

		commands.NewConnectionManager(commands.NewFS("/Users/anicolaspp")).Handle(conn)
	}
}

//func handleConnection(conn net.Conn) {
//	// close connection on exit
//	defer conn.Close()
//
//	conn.Write([]byte("220 Welcome to Nico FTP Server\n"))
//
//	cmds := make(chan []byte)
//
//	var buf [512]byte
//	for {
//		// read upto 512 bytes
//		n, err := conn.Read(buf[0:])
//
//		fmt.Println(n)
//
//		if err != nil {
//			fmt.Println(err)
//		}
//
//		cmd := string(buf[0:n])
//
//		fmt.Println(cmd)
//
//		go commands.ProcessCommand(buf[0:n], cmds)
//
//		reponse := <-cmds
//
//		fmt.Println("Sending: " + string(reponse))
//
//		conn.Write(reponse)
//
//		// write the n bytes read
//		_, err2 := conn.Write(buf[0:n])
//		if err2 != nil {
//			return
//		}
//	}
//}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
	}
}
