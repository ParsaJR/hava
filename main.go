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
		log.Fatal(".env cant be loaded. Is this file exist?")
	}
	if *help || len(os.Args) <= 1 {
		showHelp()
		return
	} else {
		apikey := os.Getenv("KEY")
		weatherapi.ReturnTemperatureByCity(os.Args[1], apikey)
		fmt.Println("󰔄")
	}
}
func showHelp() {
	fmt.Println("Usage: hava [option] [city]")
	flag.PrintDefaults()
}
