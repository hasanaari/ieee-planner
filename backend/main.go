package main

import (
	"fmt"
	"log"
	"net/http"

	// "os"
	"time"

	// "github.com/nynniaw12/ieee-planner/cache"
	"github.com/joho/godotenv" // package for loading .env
	"github.com/nynniaw12/ieee-planner/api/handlers"
	"github.com/nynniaw12/ieee-planner/db"

	"github.com/nynniaw12/ieee-planner/scraper"

	// _ "github.com/lib/pq"
	"github.com/nynniaw12/ieee-planner/middleware"
)

func StartDaemon(timeout time.Duration, f func() error) {
	go func() {
		for {
			err := f()
			if err != nil {
				log.Printf("Error executing daemon function: %v", err)
			}

			time.Sleep(timeout)
		}
	}()
}

// TODO: stores are better off in the database but this is fine for the demo
// TODO: currently scrapers are just cli tools actually run them in daemons and have proper caching mechanisms
// TODO: big todo is to have a way better major requirements scraper which is very very hard
func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file", err)
	}

	database := db.ConnectToDB()
	defer database.Close()

	// err = db.CreateCoursesTable(database)

	// if err != nil {
	// 	log.Fatal("Creating course table failed", err)
	// }

	// wd, err := os.Getwd()
	// if err != nil {
	// 	log.Fatalf("Error getting working directory: %v", err)
	// }
	// fmt.Printf("Working directory: %s\n", wd)

	// New feature in go 1.22, it actually handles restful APIs without needing to install dependencies
	courses_store, err := scraper.NewCoursesStore("./scraper-out/courses/")
	if err != nil {
		log.Fatalf("Error  creating courses store: %v", err)

	}
	// majorreqs_store, err := scraper.NewMajorRequirementsStore("./scraper-out/majorreqs/")
	// if err != nil {
	// 	log.Fatalf("Error  creating majorreqs store store: %v", err)
	// }
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/quarters", scraper.GetAvailableQuartersHandler(courses_store))
	mux.HandleFunc("GET /api/courses", scraper.GetCoursesByQuarterHandler(courses_store))
	mux.HandleFunc("GET /api/courses/subject", scraper.GetCoursesBySubjectHandler(courses_store))
	mux.HandleFunc("GET /api/courses/key", scraper.GetCoursesByKeyHandler(courses_store))

	// mux.HandleFunc("GET /api/majors", scraper.GetAvailableMajorsHandler(majorreqs_store))
	// mux.HandleFunc("GET /api/reqs", scraper.GetMajorRequirementsHandler(majorreqs_store))

    // mux.HandleFunc("GET /api/quarters", handlers.GetAvailableQuartersHandler(database))
    // mux.HandleFunc("GET /api/courses", handlers.GetCoursesByQuarterHandler(database))
    // mux.HandleFunc("GET /api/courses/subject", handlers.GetCoursesBySubjectHandler(database))
    // mux.HandleFunc("GET /api/courses/key", handlers.GetCoursesByKeyHandler(database))

    mux.HandleFunc("GET /api/majors", handlers.GetAvailableMajorsHandler(database))
    mux.HandleFunc("GET /api/reqs", handlers.GetMajorRequirementsHandler(database))

	// // testing handlers
	// // http.HandleFunc("GET /api/courses", handlers.CoursesHandler) // request to /courses, call CoursesHandler
	// mux.HandleFunc("GET /api/schedules", handlers.SchedulesHandler)
	// mux.HandleFunc("POST /api/generate", handlers.GenerateHandler) // openAI, I added POST for now, feel free to change it

	// // Add http paths to handlers
	// mux.HandleFunc("POST /api/courses", handlers.AddCourseHandler(database))
	// mux.HandleFunc("PUT /api/courses/{id}", handlers.EditCourseHandler(database))
	// mux.HandleFunc("DELETE /api/courses/{id}", handlers.RemoveCourseHandler(database))
	// mux.HandleFunc("DELETE /api/courses", handlers.ClearCoursesHandler(database))
	// mux.HandleFunc("GET /api/courses/{id}", handlers.GetCourseHandler(database))
	// mux.HandleFunc("GET /api/courses", handlers.GetCoursesHandler(database))

	fmt.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", middleware.CorsMiddleware(mux))) // error will stop program
}
