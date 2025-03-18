package weatherapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type WeatherResponse struct {
	Current Current `json:"current"`
}

type Current struct {
	Temp_C float64 `json:"temp_c"`
}

func ReturnTemperatureByCity(city string, apikey string) {
	call(city, apikey)
}

func call(city string, apikey string) {
	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", apikey, city)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatal("Status Code is wrong")
	}

	body, err := io.ReadAll(resp.Body)

	var weatherResponse WeatherResponse

	err = json.Unmarshal(body, &weatherResponse)
	fmt.Print(weatherResponse.Current.Temp_C)

}
