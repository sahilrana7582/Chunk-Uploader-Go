package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

type Manifest struct {
	Filename string   `json:"filename`
	Chunks   []string `json:chunks`
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// Read the request line (e.g. "GET MANIFEST file.pdf\n")
	reqLine, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("❌ Failed to read from client %s: %v\n", conn.RemoteAddr(), err)
		return
	}

	reqLine = strings.TrimSpace(reqLine)
	fmt.Printf("📩 Received request: %s\n", reqLine)

	parts := strings.SplitN(reqLine, " ", 3)
	if len(parts) < 2 || strings.ToUpper(parts[0]) != "GET" {
		conn.Write([]byte("ERROR: Invalid request\n"))
		return
	}

	command := strings.ToUpper(parts[1])

	switch command {

	case "MANIFEST":
		if len(parts) != 3 {
			conn.Write([]byte("ERROR: Filename missing for MANIFEST\n"))
			return
		}

		fileName := parts[2]
		filePath := fmt.Sprintf("manifests/%s.json", fileName)

		data, err := os.ReadFile(filePath)
		if err != nil {
			conn.Write([]byte("ERROR: Unable to read manifest\n"))
			return
		}

		data = append(data, '\n')

		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("❌ Error sending manifest:", err)
		}

	case "CHUNK":
		if len(parts) != 3 {
			conn.Write([]byte("ERROR: Chunk hash missing for CHUNK\n"))
			return
		}
		chunkHash := parts[2]

		chunkPath := fmt.Sprintf("chunks/%s.chunk", chunkHash)

		file, err := os.Open(chunkPath)
		if err != nil {
			conn.Write([]byte(fmt.Sprintf("ERROR: Could not open chunk %s\n", chunkHash)))
			return
		}
		defer file.Close()

		time.Sleep(1 * time.Second)

		buffer := make([]byte, 1024*1024*5)
		for {
			n, err := file.Read(buffer)
			if err != nil {
				if err == io.EOF {
					break
				}
				conn.Write([]byte(fmt.Sprintf("ERROR: Reading chunk %s failed\n", chunkHash)))
				return
			}

			_, writeErr := conn.Write(buffer[:n])
			if writeErr != nil {
				return
			}
		}

	default:
		conn.Write([]byte("ERROR: Unknown command\n"))
	}
}
