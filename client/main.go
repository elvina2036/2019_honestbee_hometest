package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	// TODO ** handle losing connect with server
	for {
		// read input
		fmt.Print("Command: ")
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		cmd := strings.TrimSpace(string(input))
		if err != nil {
			fmt.Println(err)
			return
		}

		// send request
		conn.Write([]byte(cmd + "\n"))

		// listen for response
		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(">>>> Response: ", response)
	}
}
