package main

import (
	"fmt"
	"github.com/anicolaspp/ftp/ftp"
	"net"
	"os"
)

func main() {

	listen, err := net.Listen("tcp", ":21")

	checkError(err)

	fmt.Println("Server running at port 21...")

	defer listen.Close()

	baseDir := "/Users/nperez/ftp"

	for {
		conn, err := listen.Accept()
		checkError(err)

		go ftp.NewConnectionManager(baseDir).Handle(conn)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
	}
}
