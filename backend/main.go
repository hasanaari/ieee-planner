package main

import (
	"fmt"
	"log"
	"net/http"

	// "os"
	"time"

	// "github.com/nynniaw12/ieee-planner/api/handlers"
	// "github.com/nynniaw12/ieee-planner/cache"
	// "github.com/nynniaw12/ieee-planner/db"

	// "github.com/joho/godotenv" // package for loading .env
	// "github.com/nynniaw12/ieee-planner/scraper"

	_ "github.com/lib/pq"
	"github.com/nynniaw12/ieee-planner/middleware"
	"github.com/nynniaw12/ieee-planner/scraper"
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

func main() {

	// err := godotenv.Load()

	// if err != nil {
	// 	log.Println("Error loading .env file", err)
	// }

	// database := db.ConnectToDB()
	// defer database.Close()

	// err = db.CreateCoursesTable(database)
	// // TODO: as we integrate other tables add calls to them from a function called db.CreateTables

	// if err != nil {
	// 	log.Fatal("Creating course table failed", err)
	// }

	// wd, err := os.Getwd()
	// if err != nil {
	// 	log.Fatalf("Error getting working directory: %v", err)
	// }
	// fmt.Printf("Working directory: %s\n", wd)

	// New feature in go 1.22, it actually handles restful APIs without needing to install dependencies
	store, err := scraper.NewCoursesStore("./scraper-out/")
	if err != nil {
		log.Fatalf("Error  creating courses store: %v", err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/courses", scraper.GetCoursesByQuarterHandler(store))
	mux.HandleFunc("GET /api/quarters", scraper.GetAvailableQuartersHandler(store))

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
