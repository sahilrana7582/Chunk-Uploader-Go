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
		fmt.Printf("Error generating chunks: %v", err)
	}

	err = reconst("reconst/Main.pdf", 40)
	if err != nil {
		fmt.Printf("Error reconstructing file: %v", err)
	} else {
		fmt.Println("File reconstructed successfully!")
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
			return fmt.Errorf("error reading file : %w", err)
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

func reconst(outputPath string, totalChunks int) error {
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating output file %s: %w", outputPath, err)
	}
	defer outFile.Close()

	for i := 0; i < totalChunks; i++ {
		chunkFileName := fmt.Sprintf("%s/chunk_%d", ChunkDir, i)
		data, err := os.ReadFile(chunkFileName)
		if err != nil {
			return err
		}
		_, err = outFile.Write(data)
		if err != nil {
			return err
		}
	}

	return nil
}
