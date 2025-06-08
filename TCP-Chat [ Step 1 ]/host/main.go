package main

import (
	"fmt"
	"net"
	"sync"
)

var (
	mu      = sync.Mutex{}
	Clients = make(map[net.Conn]string)
)

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		return
	}

	fmt.Println("💬 Chat Server is Live on Port 8080")

	for {
		conn, err := lis.Accept()
		if err != nil {
			fmt.Printf("❌ Error accepting connection: %v\n", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	serverName := conn.RemoteAddr().String()

	mu.Lock()
	Clients[conn] = serverName
	mu.Unlock()

	fmt.Printf("🔵 New connection from %s\n", serverName)

	buff := make([]byte, 1024)

	for {
		n, err := conn.Read(buff)
		if err != nil {
			fmt.Printf("🔴 Connection closed or error from %s: %v\n", serverName, err)
			break
		}

		if n == 0 {
			fmt.Printf("⚠️ Empty message from %s\n", serverName)
			break
		}

		msg := string(buff[:n])
		fmt.Printf("💬 Received from %s: %s", serverName, msg)

		broadcast(fmt.Sprintf("[%s] %s", serverName, msg), conn)
	}

	mu.Lock()
	delete(Clients, conn)
	mu.Unlock()
	fmt.Printf("🔴 %s disconnected\n", serverName)
}

func broadcast(msg string, sender net.Conn) {
	mu.Lock()
	defer mu.Unlock()

	for client := range Clients {
		if client != sender {
			_, err := client.Write([]byte(msg))
			if err != nil {
				fmt.Printf("⚠️ Error sending to %s: %v\n", Clients[client], err)
				client.Close()
				delete(Clients, client)
			}
		}
	}
}
