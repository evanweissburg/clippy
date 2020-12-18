package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide host:port.")
		return
	}

	port := arguments[1]
	connection, err := net.Dial("tcp", port)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Clippy client connected!")
	defer connection.Close()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		fmt.Fprintf(connection, text+"\n")

		message, _ := bufio.NewReader(connection).ReadString('\n')
		fmt.Print("->: " + message)
		if strings.TrimSpace(string(text)) == "STOP" {
			fmt.Println("Client exiting...")
			return
		}
	}
}
