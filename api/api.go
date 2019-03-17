package api

import (
	"2019_honestbee_hometest/command"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

// RequestHandler : handle api requests with rate limit
type RequestHandler struct {
	Rate     time.Duration
	Requests chan *Request
}

// Request : api request
type Request struct {
	Cmdtype int
	Para    []string
}

// ErrMsg
const (
	ErrUnavailbleAPI   = "This API is unavailable now."
	txtMockAPIResponse = "Hi, I'm Mock API."
)

// Max Value
const (
	MaxReqPara     = 2
	MaxRequestRate = 30
)

var urls = map[int]string{
	command.CMDWeather: "http://api.openweathermap.org/data/2.5/weather?APPID=c7dbecfa201c2091e1f24ae2f635920d&units=metric",
}

// EncodeReq : encoding the struct-type request and return json
func EncodeReq(paras []string, cmdtype int) (bool, string) {
	var result string
	var succeed bool

	switch cmdtype {
	case command.CMDWeather:
		succeed, result = GetWeatherReqWithPara(urls[cmdtype], paras[0])
		break
	case command.CMDMockAPI:
		succeed = true
		result = ""
		break
	default:
		succeed = false
		result = ""
		break
	}

	return succeed, result
}

// DeocodeResp : decoding the response json into formated string
func DeocodeResp(res io.Reader, cmdtype int) string {
	var result string
	switch cmdtype {
	case command.CMDWeather:
		fixedresp := DecodeWeather(res)
		result = "[" + fixedresp.Name + "] "
		result += "Now Temperature: " + floatToString(fixedresp.Main.Temp) + " °C, "
		if len(fixedresp.Weather) != 0 {
			result += "Weather Status: " + fixedresp.Weather[0].Main + "."
		}
		break
	default:
		result = ""
		break
	}
	return result
}

// ProcessRequests : --
func (r *RequestHandler) ProcessRequests(conn *net.Conn) {
	throttle := time.Tick(r.Rate)
	for req := range r.Requests {
		<-throttle
		go r.connExternalAPI(*req, conn)
	}
}

// ConnExternalAPI : --
func (r *RequestHandler) connExternalAPI(req Request, conn *net.Conn) {
	suc, reqstr := EncodeReq(req.Para, req.Cmdtype)

	// lack of parameters
	if !suc {
		sendResponse(conn, reqstr)
		return
	}

	// for mock api
	if req.Cmdtype == command.CMDMockAPI {
		sendResponse(conn, txtMockAPIResponse)
		return
	}

	// send request
	httpreq, err := http.NewRequest("GET", reqstr, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(httpreq)

	if err != nil {
		fmt.Println(err)
		sendResponse(conn, ErrUnavailbleAPI)
		return
	}

	defer resp.Body.Close()

	// get response
	sendResponse(conn, DeocodeResp(resp.Body, req.Cmdtype))
	return
}

func sendResponse(conn *net.Conn, res string) {
	(*conn).Write([]byte(res + "\n"))
}

func floatToString(num float64) string {
	return fmt.Sprintf("%.0f", num)
}
