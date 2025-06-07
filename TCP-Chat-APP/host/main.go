package main

import (
	"fmt"
	"net"
)

func main() {

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Errorf("Error listening: %v", err)
		return
	}

	fmt.Println("Host Server is Live on Port 8080")

	for {
		conn, err := lis.Accept()
		if err != nil {
			fmt.Errorf("Error accepting connection: %v", err)
			continue
		}
		fmt.Println("New connection established.\n")
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	serverName := conn.RemoteAddr().String()
	fmt.Printf("New connection from %s\n", serverName)
	defer conn.Close()

	buff := make([]byte, 1024)

	for {
		data, err := conn.Read(buff)
		if err != nil {
			fmt.Printf("Error reading from connection: %v\n", err)
			return
		}

		if data == 0 {
			fmt.Printf("Connection closed by %s\n", serverName)
			return
		}

		fmt.Printf("Received from %s: %s\n", serverName, string(buff[:data]))
	}
}
