package main

import (
	"fmt"
	"net"
	"os"
	"sync"
)

type Manifest struct {
	Filename string   `json:"filename"`
	Chunks   []string `json:"chunks"`
}

func main() {
	wg := sync.WaitGroup{}

	manifestArr := []string{
		"Main.pdf",
		// Add more manifest filenames here if needed
	}

	// Step 1: Process manifests — connect & handle manifest requests
	for _, manifest := range manifestArr {
		conn, err := net.Dial("tcp", "localhost:9000")
		if err != nil {
			fmt.Printf("Error connecting to server: %v\n", err)
			return
		}
		wg.Add(1)
		go func(c net.Conn, m string) {
			defer c.Close()
			handleRequestManifest(c, m, &wg)
		}(conn, manifest)
	}
	wg.Wait() // wait for all manifest handling to finish

	files, err := os.ReadDir("manifests")
	if err != nil {
		fmt.Printf("Error reading manifests directory: %v\n", err)
		return
	}

	// Download all chunks for each manifest concurrently
	for _, file := range files {
		wg.Add(1)
		go func(filename string) {
			chunkHandler("localhost:9000", filename, &wg)
		}(file.Name())
	}
	wg.Wait()

	for _, file := range files {
		err := reassembleChunks(file.Name())
		if err != nil {
			fmt.Printf("Error reassembling %s: %v\n", file.Name(), err)
			return
		}
		fmt.Printf("File %s reassembled successfully\n", file.Name())
	}
}
