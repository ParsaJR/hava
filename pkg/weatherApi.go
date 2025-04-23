package weatherapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// The Successfull Response
type WeatherResponse struct {
	Current Current `json:"current"`
}

type Current struct {
	Temp_C float64 `json:"temp_c"`
}

// Error Response From WeatherApi
type WeatherResponseError struct {
	Error Werror `json:"error"`
}

type Werror struct {
	Message string `json:"message"`
}

func ReturnTempertureByCity(city string, apikey string) float64 {
	temperture := call(city, apikey)
	return temperture
}

func call(city string, apikey string) float64{
	
	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", apikey, city)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		var weatherResponseError WeatherResponseError
		err = json.Unmarshal(body, &weatherResponseError)
		log.Fatal("Error: ", weatherResponseError.Error.Message)
	}

	var weatherResponse WeatherResponse

	err = json.Unmarshal(body, &weatherResponse)
	return weatherResponse.Current.Temp_C
}
