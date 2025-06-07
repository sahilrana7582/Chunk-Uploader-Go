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
		fmt.Printf("Error connecting to server: %v", err)
		return
	}

	fmt.Println("Connected to the server on port 8080")

	defer conn.Close()

	in := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter message to send: ")
		message, err := in.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			return
		}

		if message == "exit\n" {
			fmt.Println("Exiting the chat...")
			return
		}

		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Printf("Error sending message: %v\n", err)
			return
		}

	}
}
