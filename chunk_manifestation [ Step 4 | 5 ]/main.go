package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Manifest struct {
	Filename string   `json:"filename"`
	Chunks   []string `json:"chunks"`
}

func main() {
	startTime := time.Now()

	wg := sync.WaitGroup{}
	files, err := os.ReadDir("files")
	if err != nil {
		fmt.Printf("Error while reading directory\n")
		return
	}

	for idx, entry := range files {
		if entry.IsDir() {
			continue
		}
		filePath := filepath.Join("files", entry.Name())
		wg.Add(1)
		go processFile(filePath, &wg, idx)
	}

	wg.Wait()

	elapsed := time.Since(startTime)
	fmt.Printf("\n⏱️ Total processing time: %s\n", elapsed)
}

func processFile(path string, wg *sync.WaitGroup, idx int) {
	defer wg.Done()

	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", path, err)
		return
	}
	defer file.Close()

	buff := make([]byte, 1024*1024*10) // 10MB
	var chunkHashes []string

	for {
		n, err := file.Read(buff)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Error reading file: %v\n", err)
			return
		}
		if n == 0 {
			break
		}

		chunk := buff[:n]
		chunkHash := hashChunk(chunk)
		chunkHashes = append(chunkHashes, chunkHash)

		chunkPath := "chunks/" + chunkHash + ".chunk"
		if _, err := os.Stat(chunkPath); err == nil {
			fmt.Println("Chunk Exists For: -> ", idx)
			continue
		} else if !os.IsNotExist(err) {
			fmt.Printf("Error checking chunk existence: %v\n", err)
			return
		}

		err = os.WriteFile(chunkPath, chunk, 0644)
		if err != nil {
			fmt.Printf("Error writing chunk to file: %v\n", err)
			return
		}

		fmt.Printf("Chunk %s created successfully.\n", chunkHash)
		fmt.Println("-------------------------------------------")
	}

	manifest := Manifest{
		Filename: filepath.Base(path),
		Chunks:   chunkHashes,
	}

	manifestPath := "manifests/" + manifest.Filename + ".json"

	if _, err := os.Stat("manifests"); os.IsNotExist(err) {
		os.Mkdir("manifests", 0755)
	}

	manifestData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling manifest: %v\n", err)
		return
	}

	err = os.WriteFile(manifestPath, manifestData, 0644)
	if err != nil {
		fmt.Printf("Error writing manifest file: %v\n", err)
		return
	}

	fmt.Printf("📄 Manifest saved: %s\n", manifestPath)
}

func hashChunk(chunk []byte) string {
	hash := sha256.Sum256(chunk)
	return fmt.Sprintf("%x", hash)
}
