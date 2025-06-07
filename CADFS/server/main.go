package main

import (
	"fmt"
	"net"
)

const PORT = ":9000"

func main() {
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Printf("❌ Failed to start server on port %s: %v\n", PORT, err)
		return
	}
	defer listener.Close()

	fmt.Printf("🚀 Server listening on port %s...\n", PORT)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("⚠️ Error accepting connection: %v\n", err)
			continue
		}

		fmt.Printf("📥 New client connected: %s\n", conn.RemoteAddr().String())

		go handleConnection(conn)

	}
}
