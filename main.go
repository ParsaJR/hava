// Main Cli Entry

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	weatherapi "github.com/ParsaJR/hava/pkg"
	"github.com/joho/godotenv"
)

func main() {
	help := flag.Bool("h", false, "Show help")
	flag.Parse()
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env cant be loaded")
	}
	if *help {
		showHelp()
		return
	}

	City := os.Args[1]

	apikey := os.Getenv("KEY")
	weatherapi.ReturnTemperatureByCity(City, apikey)
	fmt.Println("󰔄")
}
func showHelp() {
	fmt.Println("Usage: hava [option] [City]")
	flag.PrintDefaults()
}
