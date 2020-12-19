package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func sendFile(sourceFilename string, conn net.Conn, chunkSize int) {
	file, err := os.Open(sourceFilename)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	fileChunk := make([]byte, chunkSize)
	for {
		numBytes, err := file.Read(fileChunk)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
		conn.Write(fileChunk[:numBytes])
	}
}

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Error: use command syntax 'go run tcpClient.go <host>:<port>")
		return
	}

	port := arguments[1]
	connection, err := net.Dial("tcp", port)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer connection.Close()

	fmt.Println("Clippy client sending file!")
	sendFile("test.txt", connection, 5)
}
