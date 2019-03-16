package api

import (
	"2019_honestbee_hometest/command"
	"encoding/json"
	"fmt"
	"io"
)

// MaxReqPara represents the maximum numder of Request Parameters
const MaxReqPara = 2

// Request represents a request to an external API
type Request struct {
	URL  string
	PARA []string
}

// WeatherMain save the response back from OpenWeatherMap
type WeatherMain struct {
	coord      weathercoord
	weather    weatherweather
	base       string
	Main       Weathermain
	visibility int
	wind       weatherwind
	clouds     weatherclouds
	dt         int
	sys        weathersys
	id         int
	name       string
	cod        int
}
type weathercoord struct {
	lon string
	lat string
}
type weatherweather struct {
	id          int
	main        string
	description string
	icon        string
}

// Weathermain :::
type Weathermain struct {
	Temp     int
	pressure int
	humidity int
	temp_min int
	temp_max int
}
type weatherwind struct {
	speed int
}
type weatherclouds struct {
	all int
}
type weathersys struct {
	ttype   int
	id      int
	message int
	country string
	sunrise int
	sunset  int
}

// URLS : **
var URLS = make(map[int]string)

const (
	urlWeather = "http://api.openweathermap.org/data/2.5/weather?APPID=c7dbecfa201c2091e1f24ae2f635920d&q="
)

// EncodeReq : encoding the struct-type request and return json
func EncodeReq(req Request) []byte {
	jsondata, err := json.Marshal(req)
	if err != nil {
		fmt.Println(err)
	}

	return jsondata
}

// DecodeRes : decoding the json-type response and return a struct
// DecodeWeather : decoding weather response
func DecodeWeather(res io.Reader) WeatherMain {
	structdata := WeatherMain{}
	err := json.NewDecoder(res).Decode(&structdata)
	//err := json.Unmarshal(res, &structdata)
	if err != nil {
		fmt.Println(err)
	}
	return structdata
}

// GetURL returns the urls including apiid
func GetURL(cmdtype int) string {
	return URLS[cmdtype]
}

func main() {
	// init urls
	URLS[command.CMDWeather] = urlWeather
}
