package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/nynniaw12/ieee-planner/scraper"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file", err)
	}

	major := flag.String("major", "", "Comma-separated list of quarters to whitelist")
	outputPtr := flag.String("out", "courses.json", "Output JSON file path")
	flag.Parse()

	mr, err := scraper.GetMajorreqs(*major)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	err = scraper.WriteMajorreqsToJSON(mr, *outputPtr)
	if err != nil {
		fmt.Printf("error writing major to JSON: %v\n", err)
		os.Exit(1)
	}
}
