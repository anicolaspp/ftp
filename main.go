package main

import (
	"fmt"
	"github.com/anicolaspp/ftp/ftp"
	"log"
	"net"
)

func main() {

	listener, err := net.Listen("tcp", ":21")

	checkError(err)

	fmt.Println("Server running at port 21...")

	defer listener.Close()

	baseDir := "/Users/nperez/ftp"

	for {
		conn, err := listener.Accept()
		checkError(err)

		go ftp.NewConnectionManager(baseDir).Handle(conn)
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatalln(fmt.Sprintf("Fatal error: %s", err.Error()))
	}
}
