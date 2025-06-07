package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"net"
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

	var hash1, hash2 string
	hash1, err = madman("originals/Main.pdf")
	if err != nil {
		fmt.Printf("Error calculating hash for chunk_0: %v", err)
	}

	err = reconst("reconst/Main.pdf", 40)
	if err != nil {
		fmt.Printf("Error reconstructing file: %v", err)
	}

	hash2, err = madman("reconst/Main.pdf")
	if err != nil {
		fmt.Printf("Error calculating hash: %v", err)
	}

	if hash1 == hash2 {
		fmt.Println("Hashes match! File integrity verified.")
	} else {
		fmt.Println("Hashes do not match! File integrity compromised.")
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

func madman(filePath string) (string, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filePath, err)
		return "", err
	}

	hash := sha256.Sum256(file)
	fmt.Printf("SHA-256 hash of %s: %x\n", filePath, hash)
	return fmt.Sprintf("%x", hash), nil
}

func downloadFile(conn net.Conn) {
	file, _ := os.Create("downloaded.txt")
	defer file.Close()

	for {
		lengthBuf := make([]byte, 4)
		_, err := io.ReadFull(conn, lengthBuf)
		if err != nil {
			break // EOF or broken connection
		}
		chunkSize := binary.BigEndian.Uint32(lengthBuf)
		chunk := make([]byte, chunkSize)
		io.ReadFull(conn, chunk)
		file.Write(chunk)
	}
}
