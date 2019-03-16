package main

import (
	"2019_honestbee_hometest/api"
	"2019_honestbee_hometest/command"
	"bufio"
	"fmt"
	"net"
	"net/http"
	"strings"
)

const (
	txtErrCmds = "Unavailable Commands.\n"
)

func main() {
	// TODO ** Optional PORT
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ln.Close()

	fmt.Println("Server Start.")
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	fmt.Printf("Connecting to %s\n", conn.RemoteAddr().String())
	for {
		netData, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		//cmd := strings.TrimSpace(string(netData))
		cmd := string(netData)
		s := strings.Split(cmd, " ")
		if strings.Compare(s[0], command.CmdtypeStr[command.CMDSTOP]) == 0 {
			break
		}

		// dealing command
		res := handleCmd(cmd)
		conn.Write([]byte(res + "\n"))
	}
	conn.Close()
}

func handleCmd(cmd string) string {
	// input: weather Taipei\n
	cmd = strings.TrimSuffix(cmd, "\n")
	s := strings.Split(cmd, " ")
	cmdtype := command.GetCmdType(s[0])
	if cmdtype < 0 {
		return txtErrCmds
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

	return connExternalAPI(cmdtype, para)
}

func connExternalAPI(cmdtype int, para []string) string {
	reqstruct := api.Request{URL: api.GetURL(cmdtype), PARA: para}
	requestjson := api.EncodeReq(reqstruct)

	req, err := http.NewRequest("GET", string(requestjson), nil)
	//resp, err := http.Get(string(requestjson))
	if err != nil {
		fmt.Println(err)
		return ""
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	defer resp.Body.Close()

	//var fixedresp api.WeatherMain = api.WeatherMain(api.DecodeRes(resp.Body, cmdtype))
	fixedresp := api.DecodeWeather(resp.Body)

	fmt.Println(fixedresp.Main.Temp)
	return string(fixedresp.Main.Temp)
}
