package api

import (
	"encoding/json"
	"fmt"
	"io"
)

// WeatherMain save the response back from OpenWeatherMap
type WeatherMain struct {
	coord      weathercoord
	Weather    []Weatherweather
	base       string
	Main       Weathermain
	visibility int
	wind       weatherwind
	clouds     weatherclouds
	dt         int
	sys        weathersys
	id         int
	Name       string
	cod        int
}
type weathercoord struct {
	lon float64
	lat float64
}

// Weatherweather :::
type Weatherweather struct {
	id          int
	Main        string
	description string
	icon        string
}

// Weathermain :::
type Weathermain struct {
	Temp     float64
	pressure int
	Humidity int
	tempmin  float64
	tempmax  float64
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
	message float64
	country string
	sunrise int
	sunset  int
}

const (
	errLackCityName = "Please enter a city name."
)

// GetWeatherReqWithPara : deal with the url with parameters
func GetWeatherReqWithPara(url string, location string) (bool, string) {
	if len(location) != 0 {
		return true, url + "&q=" + location
	}
	return false, errLackCityName
}

// DecodeWeather : decoding weather response
func DecodeWeather(res io.Reader) WeatherMain {
	structdata := WeatherMain{}
	err := json.NewDecoder(res).Decode(&structdata)
	if err != nil {
		fmt.Println(err)
	}
	return structdata
}
