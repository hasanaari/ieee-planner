package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/nynniaw12/ieee-planner/backend/scraper"
)

// CoursesHandler calls the scraper function and writes the result as JSON.
func CoursesHandler(w http.ResponseWriter, r *http.Request) {
	quarter := r.URL.Query().Get("quarter")
	if quarter == "" {
		quarter = "4980" // Default quarter, e.g., Spring 2025.
	}

	courseData, err := scraper.ScrapeCourses(quarter)
	if err != nil {
		http.Error(w, "Failed to scrape course data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(courseData); err != nil {
		http.Error(w, "Failed to encode result", http.StatusInternalServerError)
	}
}
