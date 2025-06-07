package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	go func() {
		for {
			msg, _ := bufio.NewReader(os.Stdin).ReadString('\n')
			fmt.Fprintf(conn, msg)
		}
	}()

	for {
		reply, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("🔴 Disconnected from server")
			return
		}
		fmt.Print(reply)
	}
}
