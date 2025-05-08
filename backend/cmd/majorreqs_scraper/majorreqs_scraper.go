package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/nynniaw12/ieee-planner/db"
	"github.com/nynniaw12/ieee-planner/scraper"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Parse command line flags
	major := flag.String("major", "", "Major to retrieve requirements for")
	flag.Parse()

	if *major == "" {
		fmt.Println("Error: Major parameter is required")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Connect to the database
	database := db.ConnectToDB()
	defer database.Close()

	// Create the major requirements table if it doesn't exist
	err = db.CreateMajorReqsTableIfNotExists(database)
	if err != nil {
		fmt.Printf("Error creating major requirements table: %v\n", err)
		os.Exit(1)
	}

	// Get major requirements from scraper
	mr, err := scraper.GetMajorreqs(*major)
	if err != nil {
		fmt.Printf("Error getting major requirements: %v\n", err)
		os.Exit(1)
	}

	// Write the major requirements to the database
	err = db.WriteMajorReqsToDatabase(database, &mr)
	if err != nil {
		fmt.Printf("Error writing major requirements to database: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully wrote %s major requirements to database\n", *major)
}