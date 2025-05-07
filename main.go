// Main Cli Entry
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	weatherapi "github.com/ParsaJR/hava/pkg"
	"github.com/joho/godotenv"
)

func main() {
	help := flag.Bool("h", false, "Show help")
	flag.Parse()
	envCheck()
	// Return help, if user wants it or didn't provide any argument after the commnad
	if *help || len(os.Args) <= 1 {
		showHelp()
		return
	} else {
		apikey := os.Getenv("KEY")
		temp, err := weatherapi.ReturnTempertureByCity(os.Args[1], apikey)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%0.1f°C \n", temp)
	}
}

func showHelp() {
	fmt.Println("Usage: hava [city] [option]")
	flag.PrintDefaults()
}

func envCheck() {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	err = godotenv.Load()
	if err != nil {
		for range 2 {
			err = godotenv.Load(filepath.Join(ex, "/.env"))
			if err != nil {
				ex = filepath.Join(ex, "../")
			} else {
				break
			}
		}
	}
	if err != nil {
		log.Fatal(err)
	}
}
