package main

import (
	"2019_honestbee_hometest/api"
	"2019_honestbee_hometest/command"
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

const (
	txtErrCmds = "Unavailable Commands.\n"
)

var connectioncnt int
var processedreqcnt int

// limiter
var requestsCh = make(chan *api.Request, 1000)
var rate = time.Second / api.MaxRequestRate // handle 30 tx per second
var handler = &api.RequestHandler{Rate: rate, Requests: requestsCh}

func main() {
	handleHTTPListener()
	handleTCPListener()
}

func handleTCPListener() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ln.Close()

	fmt.Println("TCP Server On.")
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleTCPConnection(conn)
	}
}

func handleTCPConnection(conn net.Conn) {
	fmt.Printf("Connecting to %s\n", conn.RemoteAddr().String())
	connectioncnt++
	for {
		netData, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			break
		}

		cmd := strings.TrimSuffix(string(netData), "\n")
		s := strings.Split(cmd, " ")
		if strings.Compare(s[0], command.CmdtypeStr[command.CMDSTOP]) == 0 {
			break
		}

		// dealing command
		go handleCmd(cmd, &conn)
	}
	fmt.Printf("%s disconnect.\n", conn.RemoteAddr().String())
	conn.Close()
	connectioncnt--
}

func handleCmd(cmd string, conn *net.Conn) {
	// input: weather Taipei\n
	cmd = strings.TrimSuffix(cmd, "\n")
	s := strings.Split(cmd, " ")
	cmdtype := command.GetCmdType(s[0])
	if cmdtype < 0 {
		(*conn).Write([]byte(txtErrCmds + "\n"))
		return
	}
	var para []string
	slen := len(s)
	for i := 1; i <= api.MaxReqPara; i++ {
		if i < slen {
			para = append(para, s[i])
		} else {
			para = append(para, "")
		}
	}

	handler.Requests <- &api.Request{Cmdtype: cmdtype, Para: para}
	processedreqcnt++
	go handler.ProcessRequests(conn)
	return
}
