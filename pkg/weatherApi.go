package weatherapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
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

func ReturnTempertureByCity(city string, apikey string) (float64, error) {
	cacheFile := "weather_cache"
	var err error

	if !isCacheValid(cacheFile) {
		err = dump(city, apikey)
		if err != nil {
			return 0, err
		}
	}

	file, err := os.ReadFile(cacheFile)
	if err != nil {
		return 0, err
	}

	var weather WeatherResponse
	err = json.Unmarshal(file, &weather)

	return weather.Current.Temp_C, nil

}

func dump(city string, apikey string) error {

	// this function will Fetch The response and write it into the disk...

	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", apikey, city)

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		// trying to parse the error
		var apierror WeatherResponseError
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("failed to read the response: %w", err)
		}
		if err = json.Unmarshal(body, &apierror); err != nil {
			return fmt.Errorf("Failed to parse Api response: %w",err)
		}
		return fmt.Errorf("Api Error: %s", apierror.Error.Message)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = os.WriteFile("weather_cache", body, 0644)
	if err != nil {
		return err
	}

	return nil
}

func isCacheValid(cachefile string) bool {
	if _, err := os.Stat(cachefile); os.IsNotExist(err) {
		return false
	}

	fileinfo, _ := os.Stat(cachefile)
	return time.Since(fileinfo.ModTime()) < 3*time.Minute
}
