/* package main

import (
	"fmt"
    "log"
    "net/http"
    "os"
	
	"github.com/joho/godotenv" // package for loading .env
	"github.com/nynniaw12/ieee-planner/backend/api/handlers"
)

func main(){
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// testing handlers
	http.HandleFunc("/api/courses", handlers.CoursesHandler) // request to /courses, call CoursesHandler
    http.HandleFunc("/api/schedules", handlers.SchedulesHandler)
	http.HandleFunc("/api/generate", GenerateHandler) // openAI


    fmt.Println("Server starting on :8080...")
    log.Fatal(http.ListenAndServe(":8080", nil)) // error will stop program

} */

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nynniaw12/ieee-planner/backend/scraper"
	// other imports, such as your handlers for other routes
)

func main() {
	// Register your course handler for testing purposes.
	http.HandleFunc("/test-courses", scraper.CourseHandler)

	// If you already have other routes, you can keep them here as well.
	// http.HandleFunc("/api/courses", yourExistingCoursesHandler)

	fmt.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
