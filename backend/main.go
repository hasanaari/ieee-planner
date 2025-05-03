package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nynniaw12/ieee-planner/api/handlers"
	"github.com/nynniaw12/ieee-planner/cache"
	"github.com/nynniaw12/ieee-planner/db"

	"github.com/gocolly/colly/v2"
	"github.com/joho/godotenv" // package for loading .env
	"github.com/nynniaw12/ieee-planner/scraper"

	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Println("Error loading .env file", err)
	}


	database := db.ConnectToDB()
	defer database.Close()

	err = db.CreateCoursesTable(database)

	if err != nil {
		log.Fatal("Creating course table failed", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting working directory: %v", err)
	}
	fmt.Printf("Working directory: %s\n", wd)

	// Instrument Colly to log every URL it visits:
	scraper.SetLogger(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL.String())
	})

	log.Println("Starting full scrape of courses (may take a couple minutes)…")
	scraper.CachedCourses, scraper.ScrapeError = scraper.ScrapeCourses("4980")
	if scraper.ScrapeError != nil {
		log.Fatalf("Initial scrape failed: %v", scraper.ScrapeError)
	}
	log.Printf("Scrape complete: %d courses loaded\n", len(scraper.CachedCourses))

	// New feature in go 1.22, it actually handles restful APIs without needing to install dependencies
	mux := http.NewServeMux()	

	// testing handlers
	// http.HandleFunc("GET /api/courses", handlers.CoursesHandler) // request to /courses, call CoursesHandler
	mux.HandleFunc("GET /api/schedules", handlers.SchedulesHandler)
	mux.HandleFunc("POST /api/generate", handlers.GenerateHandler) // openAI, I added POST for now, feel free to change it

	// Add http paths to handlers
	mux.HandleFunc("POST /api/courses", handlers.AddCourseHandler(database))
	mux.HandleFunc("PUT /api/courses/{id}", handlers.EditCourseHandler(database))
	mux.HandleFunc("DELETE /api/courses/{id}", handlers.RemoveCourseHandler(database))
	mux.HandleFunc("DELETE /api/courses", handlers.ClearCoursesHandler(database))
	mux.HandleFunc("GET /api/courses/{id}", handlers.GetCourseHandler(database))
	mux.HandleFunc("GET /api/courses", handlers.GetCoursesHandler(database))
	
	if !cache.IsCacheValid(database, "courses"){
		log.Println("Invalid cached courses, rescraping...")
		scraper.CachedCourses, scraper.ScrapeError = scraper.ScrapeCourses("4980")
		if scraper.ScrapeError != nil {
			log.Fatalf("Rescraping failed: %v", scraper.ScrapeError)
		}
		log.Printf("Rescraping complete: %d courses loaded\n", len(scraper.CachedCourses))
	}

	if !cache.IsCacheValid(database, "MajorRequirements"){
		// For when scraping major requirements is done
	}

	fmt.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", mux)) // error will stop program

}


