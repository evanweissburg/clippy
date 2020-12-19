package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

var connectionCount = 0

func receiveFile(conn net.Conn, targetFilename string, chunkSize int) {
	file, err := os.Create(targetFilename)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
	}

	fileChunk := make([]byte, chunkSize)
	for {
		numBytes, err := conn.Read(fileChunk)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
		file.Write(fileChunk[:numBytes])
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Printf("[Server] Received file %d\n", connectionCount)
	receiveFile(conn, "recieved.txt", 5)
}

func main() {
	arguments := os.Args
	if len(arguments) != 2 {
		fmt.Println("Error: use command syntax 'go run tcpServer.go <port>'")
		return
	}
	port := ":" + arguments[1]

	listener, err := net.Listen("tcp4", port)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Clippy server online!")
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		connectionCount++
		go handleConnection(conn)
	}
}
