package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nynniaw12/ieee-planner/api/handlers"
	"github.com/nynniaw12/ieee-planner/db"

	"github.com/gocolly/colly/v2"
	"github.com/joho/godotenv" // package for loading .env
	"github.com/nynniaw12/ieee-planner/scraper"

	_ "github.com/lib/pq"
)

func main() {
	database := db.ConnectToDB()
	defer database.Close()

	err := db.CreateCoursesTable(database)

	if err != nil {
		log.Fatal("Creating course table failed", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting working directory: %v", err)
	}
	fmt.Printf("Working directory: %s\n", wd)

	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

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

	// testing handlers
	http.HandleFunc("/api/courses", handlers.CoursesHandler) // request to /courses, call CoursesHandler
	http.HandleFunc("/api/schedules", handlers.SchedulesHandler)
	http.HandleFunc("/api/generate", handlers.GenerateHandler) // openAI

	// Add http paths to handlers
	http.HandleFunc("/add-course", handlers.AddCourseHandler(database))
	http.HandleFunc("/edit-course", handlers.EditCourseHandler(database))
	http.HandleFunc("/remove-course", handlers.RemoveCourseHandler(database))
	http.HandleFunc("/clear-courses", handlers.ClearCoursesHandler(database))
	http.HandleFunc("/get-course", handlers.GetCourseHandler(database))
	http.HandleFunc("/get-courses", handlers.GetCoursesHandler(database))
	
	fmt.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil)) // error will stop program

}


