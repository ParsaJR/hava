package weatherapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestApi(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Print("no .env file exist")
	}
	apikey := os.Getenv("KEY")
	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", apikey, "Tehran")

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected %d status code, got %d",http.StatusOK,resp.StatusCode)	
	}
	var weatherResponse WeatherResponse

	if err := json.Unmarshal(body, &weatherResponse); err != nil {
		t.Fatalf("Failed to unmarshal response: %v",err)	
	} 

}
