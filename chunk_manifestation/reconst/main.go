package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Manifest struct {
	Filename string   `json:"filename"`
	Chunks   []string `json:"chunks"`
}

func main() {
	manifestPath := filepath.Join("../manifests", "Main.pdf.json")
	manifestFile, err := os.Open(manifestPath)
	if err != nil {
		fmt.Printf("❌ Error opening manifest: %v\n", err)
		return
	}
	defer manifestFile.Close()

	var manifest Manifest
	if err := json.NewDecoder(manifestFile).Decode(&manifest); err != nil {
		fmt.Printf("❌ Error parsing manifest: %v\n", err)
		return
	}

	outputPath := filepath.Join("reconstructed", manifest.Filename)
	if err := os.MkdirAll("reconstructed", 0755); err != nil {
		fmt.Printf("❌ Error creating output directory: %v\n", err)
		return
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("❌ Error creating output file: %v\n", err)
		return
	}
	defer outputFile.Close()

	for idx, chunkHash := range manifest.Chunks {
		chunkPath := filepath.Join("../chunks", chunkHash+".chunk")
		chunkFile, err := os.Open(chunkPath)
		if err != nil {
			fmt.Printf("❌ Error opening chunk %d (%s): %v\n", idx, chunkHash, err)
			return
		}

		_, err = io.Copy(outputFile, chunkFile)
		chunkFile.Close()
		if err != nil {
			fmt.Printf("❌ Error writing chunk %d: %v\n", idx, err)
			return
		}

		fmt.Printf("✅ Chunk %d (%s) appended successfully\n", idx+1, chunkHash)
	}

	fmt.Printf("\n🎉 Reconstructed file written to: %s\n", outputPath)
}
