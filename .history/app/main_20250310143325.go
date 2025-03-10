package main

import (
	"fmt"
	"net"
	"os"
)

func main() {

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	// Send HTTP 200 OK response
	response := "HTTP/1.1 200 OK\r\n\r\n"
	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
		os.Exit(1)
	}

	// Close the connection
	conn.Close()
}
