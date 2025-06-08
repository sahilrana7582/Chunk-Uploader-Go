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

func chunkHandler(serverAddr, filepath string, wg *sync.WaitGroup) {
	defer wg.Done()

	fileName := fmt.Sprintf("manifests/%s", filepath)
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("❌ Error opening manifest: %v\n", err)
		return
	}
	defer file.Close()

	var manifest Manifest
	if err := json.NewDecoder(file).Decode(&manifest); err != nil {
		fmt.Printf("❌ Error decoding manifest: %v\n", err)
		return
	}

	// Create download directory if not exists
	_ = os.MkdirAll("downloads", os.ModePerm)

	var chunkWG sync.WaitGroup
	for _, chunk := range manifest.Chunks {
		chunkWG.Add(1)
		go func(chunk string) {
			defer chunkWG.Done()

			data, err := dataReceive(serverAddr, chunk)
			if err != nil {
				log.Printf("❌ Error fetching chunk %s: %v\n", chunk, err)
				return
			}

			err = os.WriteFile(fmt.Sprintf("downloads/%s.chunk", chunk), data, 0644)
			if err != nil {
				log.Printf("❌ Error writing chunk %s: %v\n", chunk, err)
			} else {
				fmt.Printf("✅ Chunk %s saved.\n", chunk)
			}
		}(chunk)
	}

	chunkWG.Wait() // wait for all chunks to finish downloading before returning
}

func dataReceive(serverAddr, chunk string) ([]byte, error) {
	localPath := fmt.Sprintf("downloads/%s.chunk", chunk)
	if file, err := os.Open(localPath); err == nil {
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			return nil, fmt.Errorf("❌ Error reading local chunk %s: %v", chunk, err)
		}

		fmt.Printf("📁 Chunk %s read from local disk\n", chunk)
		return data, nil
	}

	fmt.Printf("🌐 Chunk %s not found locally — requesting from server\n", chunk)

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return nil, fmt.Errorf("❌ Could not connect to server: %v", err)
	}
	defer conn.Close()

	req := fmt.Sprintf("GET CHUNK %s\n", chunk)
	_, err = conn.Write([]byte(req))
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to send request: %v", err)
	}

	var buffer []byte
	temp := make([]byte, 1024*1024)
	for {
		n, err := conn.Read(temp)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("❌ Error receiving chunk %s: %v", chunk, err)
		}
		buffer = append(buffer, temp[:n]...)
	}

	return buffer, nil
}

func reassembleChunks(filepath string) error {
	manifestPath := fmt.Sprintf("manifests/%s", filepath)

	file, err := os.Open(manifestPath)
	if err != nil {
		return fmt.Errorf("could not open manifest: %v", err)
	}
	defer file.Close()

	var manifest Manifest
	if err := json.NewDecoder(file).Decode(&manifest); err != nil {
		return fmt.Errorf("could not decode manifest: %v", err)
	}

	_ = os.MkdirAll("files", os.ModePerm)
	finalPath := fmt.Sprintf("files/%s", manifest.Filename)
	out, err := os.Create(finalPath)
	if err != nil {
		return fmt.Errorf("could not create final file: %v", err)
	}
	defer out.Close()

	for _, chunk := range manifest.Chunks {
		data, err := os.ReadFile(fmt.Sprintf("downloads/%s.chunk", chunk))
		if err != nil {
			return fmt.Errorf("error reading chunk %s: %v", chunk, err)
		}
		_, err = out.Write(data)
		if err != nil {
			return fmt.Errorf("error writing chunk %s: %v", chunk, err)
		}
	}

	return nil
}
