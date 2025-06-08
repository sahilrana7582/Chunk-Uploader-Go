package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

func main() {
	for i := 0; i < 10; i++ {
		moreServer(i)
	}
	fmt.Println("File download completed.")
}

func moreServer(n int) {

	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		fmt.Printf("Error connecting to server: %v\n", err)
		return
	}

	defer conn.Close()

	conn.Write([]byte("GET Main.pdf\n"))

	downloadFile(conn, n)
}
func downloadFile(conn net.Conn, n int) {
	fmt.Printf("Downloading file from server %d...\n", n)
	if _, err := os.Stat("dowloads"); os.IsNotExist(err) {
		err := os.Mkdir("dowloads", 0755)
		if err != nil {
			fmt.Printf("Error creating directory: %v\n", err)
			return
		}
	}
	file, _ := os.Create("dowloads/downloaded" + strconv.Itoa(n) + ".pdf")
	defer file.Close()

	buff := make([]byte, 1024*1024*10) // 10 MB buffer

	for {

		n, err := conn.Read(buff)
		if err != nil {
			if err.Error() != "EOF" {
				fmt.Printf("Error reading from connection: %v\n", err)
			}
			break
		}
		if n == 0 {
			fmt.Println("No more data to read.")
			break
		}
		if n < len(buff) {
			buff = buff[:n]
		}
		_, err = file.Write(buff)
		if err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
			return
		}
		fmt.Printf("Wrote %d bytes to file\n", n)
	}
}
