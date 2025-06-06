package main

import (
	"fmt"
	"io"
	"os"
)

const (
	ChunkSize      = 1 * 1024 * 1024 * 4
	ChunkDir       = "chunks"
	ReconstructDir = "reconst"
)

func main() {
	err := generateChunks("originals/Main.pdf")
	if err != nil {
		fmt.Println("Error generating chunks: %v", err)
		return
	}
}

func generateChunks(filepath string) error {

	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	totalPartsNum := int((fileInfo.Size() + ChunkSize - 1) / ChunkSize)
	fmt.Printf("Splitting into %d chunks\n", totalPartsNum)

	buff := make([]byte, ChunkSize)

	var index int

	for {
		byteRead, err := file.Read(buff)
		if err != nil && err != io.EOF {
			return fmt.Errorf("error reading file: %w", err)
		}

		if byteRead == 0 {
			break
		}

		chunkFileName := fmt.Sprintf("%s/chunk_%d", ChunkDir, index)

		err = os.WriteFile(chunkFileName, buff[:byteRead], 0644)
		if err != nil {
			return fmt.Errorf("error writing chunk file %s: %w", chunkFileName, err)
		}
		index++
	}

	return nil
}
