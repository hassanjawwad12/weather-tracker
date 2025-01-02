package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type apiConfigData struct {
	OpenWeatherApiKey string `json:"OpenWeatherApiKey"`
}

type WeatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

func loadApiConfig(filename string) (apiConfigData, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return apiConfigData{}, err
	}

	var c apiConfigData

	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return apiConfigData{}, err
	}
	return c, nil
}

func GetCity(w http.ResponseWriter, r *http.Request) {
	city := strings.SplitN(r.URL.Path, "/", 3)[2]
	data, err := query(city)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)

}

func query(city string) (WeatherData, error) {

	apiConfig, err := loadApiConfig(".apiConfig")
	if err != nil {
		return WeatherData{}, err
	}

	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + apiConfig.OpenWeatherApiKey + "&q=" + city)

	if err != nil {
		return WeatherData{}, err
	}

	defer resp.Body.Close()

	var d WeatherData
	if json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return WeatherData{}, err
	}
	return d, nil
}

func main() {
	fmt.Println("Weather Tracker Project")
	http.HandleFunc("/weather/", GetCity)
	http.ListenAndServe(":8080", nil)
}
