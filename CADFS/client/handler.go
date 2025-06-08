package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

func handleRequestManifest(conn net.Conn, filename string, wg *sync.WaitGroup) {
	defer conn.Close()
	defer wg.Done()

	reqBody := fmt.Sprintf("GET MANIFEST %s\n", filename)

	_, err := conn.Write([]byte(reqBody))

	if err != nil {
		fmt.Printf("Error While Requesting Manifest: %v", err)
		return
	}

	data, err := io.ReadAll(conn)
	if err != nil {
		log.Fatalf("❌ Error reading from server: %v", err)
	}

	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		log.Fatalf("❌ Failed to parse JSON: %v", err)
	}

	filePath := fmt.Sprintf("manifests/%s.json", filename)
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Not Able To Create The File")
		return
	}

	defer file.Close()

	jsonData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		fmt.Printf("❌ Failed to marshal manifest: %v\n", err)
		return
	}

	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Printf("❌ Failed to write JSON to file: %v\n", err)
		return
	}

	fmt.Println("✅ Manifest saved to file successfully.")

	time.Sleep(4 * time.Second)
}
