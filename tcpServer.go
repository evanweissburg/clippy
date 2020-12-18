// https://www.linode.com/docs/guides/developing-udp-and-tcp-clients-and-servers-in-go/

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

var connectionCount = 0

func handleConnection(c net.Conn) {
	fmt.Printf("[Server] Received connection %d\n", connectionCount)
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		temp := strings.TrimSpace(string(netData))
		if temp == "STOP" {
			break
		}
		fmt.Println(temp)
		counter := strconv.Itoa(connectionCount) + "\n"
		c.Write([]byte(string(counter)))
	}
	c.Close()
}

func main() {
	arguments := os.Args
	if len(arguments) != 2 {
		fmt.Println("Error: use command syntax 'go run udpServer.go port'")
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
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		connectionCount++
		go handleConnection(connection)
	}
}
