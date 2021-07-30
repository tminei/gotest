package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Weather struct {
	Current [1]CurrentCondition `json:"current_condition"`
}

type CurrentCondition struct {
	Temperature string `json:"temp_C"`
	Visibility  string `json:"visibility"`
	Humidity    string `json:"humidity"`
	Pressure    string `json:"pressure"`
}

var lastWeather Weather

func updateWeatherData(city string) {
	for {
		response, err := http.Get(fmt.Sprintf("https://www.wttr.in/%s?format=j1", city))
		if err != nil {
			log.Fatal(err)
		}
		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		json.Unmarshal(responseData, &lastWeather)
		time.Sleep(5 * time.Minute)
	}
}

func httpServer(city string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, fmt.Sprintf(
			"City: %s; Temperature: %s; Visibility: %s; Humidity: %s; Pressure: %s;",
			city,
			lastWeather.Current[0].Temperature,
			lastWeather.Current[0].Visibility,
			lastWeather.Current[0].Humidity,
			lastWeather.Current[0].Pressure))
	})
	http.ListenAndServe(":1488", nil)
}

func main() {
	var city string = "Боярка"
	go updateWeatherData(city)
	httpServer(city)
}
