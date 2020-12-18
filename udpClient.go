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
	if len(arguments) != 2 {
		fmt.Println("Error: use command syntax 'go run udpClient.go host:port'")
		return
	}
	address := arguments[1]

	resolvedAddress, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		fmt.Println(err)
		return
	}

	connection, err := net.DialUDP("udp4", nil, resolvedAddress)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("The UDP server is %s\n", connection.RemoteAddr().String())
	defer connection.Close()

	for {
		inputReader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := inputReader.ReadString('\n')

		data := []byte(text + "\n")
		_, err = connection.Write(data)
		if strings.TrimSpace(string(data)) == "STOP" {
			fmt.Println("Exiting UDP client!")
			return
		}
		if err != nil {
			fmt.Println(err)
			return
		}

		buffer := make([]byte, 1024)
		n, _, err := connection.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Reply: %s\n", string(buffer[0:n]))
	}
}
