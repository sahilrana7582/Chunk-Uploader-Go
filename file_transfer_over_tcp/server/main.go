package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	listher, err := net.Listen("tcp", ":9000")
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		return
	}

	defer listher.Close()
	fmt.Println("Server is listening on port 9000...")

	for {
		conn, err := listher.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}

		fmt.Printf("Client connected: %s\n", conn.RemoteAddr().String())

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	req, _ := reader.ReadString('\n')
	req = strings.TrimSpace(req)

	if strings.HasPrefix(req, "GET ") {
		fileName := strings.TrimPrefix(req, "GET ")
		fileName = strings.TrimSuffix(fileName, "\n")
		fmt.Printf("Requested file: %s\n", fileName)
		file, err := os.Open("files/" + fileName)

		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			conn.Write([]byte("ERROR: File not found\n"))
			return
		}

		defer file.Close()

		buff := make([]byte, 1024*1024*10)

		for {
			n, err := file.Read(buff)
			if err != nil {
				if err.Error() != "EOF" {
					fmt.Printf("Error reading file: %v\n", err)
				}
				break
			}

			if n == 0 {
				break
			}

			// Send the file content to the client
			_, err = conn.Write(buff[:n])
			if err != nil {
				fmt.Printf("Error sending file: %v\n", err)
				return
			}
			fmt.Printf("Sent %d bytes of file %s to client\n", n, fileName)
			fmt.Println("Waiting for 2 seconds before sending more data...")
			time.Sleep(2 * time.Second)
		}

	}

}
