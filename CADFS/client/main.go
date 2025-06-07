package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type Manifest struct {
	Filename string   `json:"filename`
	Chunks   []string `json:chunks`
}

func main() {

	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	n, err := conn.Write([]byte("GET MANIFEST Main.pdf\n"))
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	fmt.Println("Number of bytes Wrote: ", n)
	time.Sleep(3 * time.Second)

	for {
		reader := bufio.NewReader(conn)

		// Read the request line (e.g. "GET MANIFEST file.pdf\n")
		reqLine, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Printf("❌ Failed to read from client %s: %v\n", conn.RemoteAddr(), err)
			return
		}

		var metadata Manifest

		err = json.Unmarshal(reqLine, &metadata)
		if err != nil {
			fmt.Println("JSON unmarshal error:", err)
			return
		}

		fmt.Printf("Data: %v", metadata)

		time.Sleep(4 * time.Second)
	}

}
