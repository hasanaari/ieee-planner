package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/nynniaw12/ieee-planner/db"
	"github.com/nynniaw12/ieee-planner/scraper"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file", err)
	}

	database := db.ConnectToDB()
	defer database.Close()

	quartersPtr := flag.String("quarters", "", "Comma-separated list of quarters to whitelist")
	schoolsPtr := flag.String("schools", "", "Comma-separated list of schools to whitelist")
	// outputPtr := flag.String("out", "courses.json", "Output JSON file path") // Not needed anymore, writing to DB

	flag.Parse()

	var whitelistedQuarters []string
	if *quartersPtr != "" {
		whitelistedQuarters = strings.Split(*quartersPtr, ",")
		for i, q := range whitelistedQuarters {
			whitelistedQuarters[i] = strings.TrimSpace(q)
		}
	}

	var whitelistedSchools []string
	if *schoolsPtr != "" {
		whitelistedSchools = strings.Split(*schoolsPtr, ",")
		for i, s := range whitelistedSchools {
			whitelistedSchools[i] = strings.TrimSpace(s)
		}
	}

	courses := scraper.ScrapeCourseDescriptionHierarchy(whitelistedQuarters, whitelistedSchools)
	for _, course := range courses {
		err := db.WriteCourseDataToDatabase(database, *course)

		if err != nil {
			fmt.Printf("error writing courses to database: %v\n", err)
			os.Exit(1)
		}
	}
	// err := scraper.WriteCourseDataToJSON(courses, *outputPtr)
	// if err != nil {
	// 	fmt.Printf("error writing courses to JSON: %v\n", err)
	// 	os.Exit(1)
	// }
}
