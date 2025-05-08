package weatherapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// The Successfull Response
type WeatherResponse struct {
	Current  Current  `json:"current"`
	Location Location `json:"location"`
}

type Location struct {
	Tz_ID   string `json:"tz_id"`
	Name    string `json:"name"`
	Country string `json:"country"`
}

type Current struct {
	Temp_C float64 `json:"temp_c"`
	Temp_F float64 `json:"temp_f"`
}

// Error Response From WeatherApi
type WeatherResponseError struct {
	Error Werror `json:"error"`
}

type Werror struct {
	Message string `json:"message"`
}

type CachedWeatherData struct {
	RequestedCity string          `json:"requested_city"`
	Weather       WeatherResponse `json:"weather"`
}

func ReturnTemperatureByCity(city string, apikey string) (float64, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir, err = os.Executable()
		if err != nil {
			return 0, fmt.Errorf("Cannot determine the executable path: %w", err)
		}
	}
	cacheFile := filepath.Join(cacheDir, "weather_cache")

	if !isCacheValid(cacheFile, city) {
		err = dump(city, apikey, cacheFile)
		if err != nil {
			return 0, err
		}
	}

	file, err := os.ReadFile(cacheFile)
	if err != nil {
		return 0, err
	}

	var cachedData CachedWeatherData
	err = json.Unmarshal(file, &cachedData)

	if err != nil {
		return 0, fmt.Errorf("failed to unmarshal weather data: %w", err)
	}
	return cachedData.Weather.Current.Temp_C, nil

}

func dump(city string, apikey string, cachefile string) error {

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
			return fmt.Errorf("Failed to parse Api response: %w", err)
		}
		return fmt.Errorf("Api Error: %s", apierror.Error.Message)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var weather WeatherResponse
	if err = json.Unmarshal(body, &weather); err != nil {
		return fmt.Errorf("failed to unmarshal weather data: %w", err)
	}
	cache := CachedWeatherData{
		RequestedCity: city,
		Weather:       weather,
	}
	data, err := json.Marshal(cache)
	err = os.WriteFile(cachefile, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func isCacheValid(cachefile string, city string) bool {
	if _, err := os.Stat(cachefile); os.IsNotExist(err) {
		return false
	}

	cache, err := os.ReadFile(cachefile)
	if err != nil {
		return false
	}

	var data CachedWeatherData
	if err = json.Unmarshal(cache, &data); err != nil {
		return false
	}

	if strings.ToLower(data.RequestedCity) != strings.ToLower(city) {
		return false
	}
	fileinfo, _ := os.Stat(cachefile)
	return time.Since(fileinfo.ModTime()) < 3*time.Minute
}
