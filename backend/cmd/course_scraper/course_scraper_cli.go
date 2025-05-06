package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/nynniaw12/ieee-planner/scraper"
)

func main() {
	quartersPtr := flag.String("quarters", "", "Comma-separated list of quarters to whitelist")
	schoolsPtr := flag.String("schools", "", "Comma-separated list of schools to whitelist")
	outputPtr := flag.String("out", "courses.json", "Output JSON file path")

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
	err := scraper.WriteCourseDataToJSON(courses, *outputPtr)
	if err != nil {
		fmt.Printf("error writing courses to JSON: %v\n", err)
		os.Exit(1)
	}
}
