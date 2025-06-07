package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
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

		data := Manifest{
			Filename: "Main.pdf",
			Chunks: []string{
				"Chunk1",
				"Chunk1",
				"Chunk1",
				"Chunk1",
				"Chunk1",
			},
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Println("JSON marshal error:", err)
			return
		}

		jsonData = append(jsonData, '\n')

		for {
			_, err = conn.Write(jsonData)
			if err != nil {
				fmt.Println("Write error:", err)
				return
			}

			fmt.Println("Sent message:", string(jsonData))
			time.Sleep(4 * time.Second)
		}

		// fileName := parts[2]
		// sendManifest(conn, fileName)

	case "CHUNK":
		if len(parts) != 3 {
			conn.Write([]byte("ERROR: Chunk hash missing for CHUNK\n"))
			return
		}

		for {
			conn.Write([]byte("Chunk Data.....\n"))
			time.Sleep(2 * time.Second)
		}
		// chunkHash := parts[2]
		// sendChunk(conn, chunkHash)

	default:
		conn.Write([]byte("ERROR: Unknown command\n"))
	}
}
