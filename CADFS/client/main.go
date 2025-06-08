package main

import (
	"fmt"
	"net"
	"sync"
)

type Manifest struct {
	Filename string   `json:"filename`
	Chunks   []string `json:chunks`
}

func main() {

	wg := sync.WaitGroup{}

	menifestArr := []string{
		"Main.pdf",
		"one.pdf",
		"two.pdf",
		"three.pdf",
		"four.pdf",
		"five.pdf",
	}

	for _, manifest := range menifestArr {
		conn, err := net.Dial("tcp", "localhost:9000")
		if err != nil {
			fmt.Printf("Error: %v", err)
			return
		}
		wg.Add(1)
		go handleRequestManifest(conn, manifest, &wg)
	}

	wg.Wait()

}
