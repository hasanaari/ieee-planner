package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/nynniaw12/ieee-planner/api/models"
)

// Add course
func AddCourseHandler(db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		
		decoder := json.NewDecoder(r.Body)

		var course models.Course

		err := decoder.Decode(&course)

		if err != nil {
			http.Error(w, "Valid course not entered", http.StatusBadRequest)
			return
		}

		err = models.AddCourse(db, course)

		if err != nil {
			http.Error(w, "Adding course failed", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Course added successfully"})
	}
}

func EditCourseHandler(db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		if r.Method != http.MethodPut{
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		decoder := json.NewDecoder(r.Body)

		var course models.Course

		err := decoder.Decode(&course)

		if err != nil {
			http.Error(w, "Valid course not entered", http.StatusBadRequest)
			return
		}

		err = models.EditCourse(db, course)

		if err != nil {
			http.Error(w, "Edit course failed", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]string{"message": "Course edited successfully"})
	}
}

func RemoveCourseHandler(db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		if r.Method != http.MethodDelete{
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		decoder := json.NewDecoder(r.Body)

		var id int

		err := decoder.Decode(&id)

		if err != nil {
			http.Error(w, "Not valid id", http.StatusBadRequest)
			return
		}

		err = models.RemoveCourse(db, id)

		if err != nil {
			http.Error(w, "Remove course failed", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Course removed successfully"})
	}
}

func ClearCoursesHandler(db *sql.DB)http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		if r.Method != http.MethodDelete{
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		err := models.ClearCourses(db)

		if err != nil {
			http.Error(w, "Clear courses failed", http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Courses cleared successfully"})
	}
}

func GetCourseHandler(db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		if r.Method != http.MethodGet{
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var id int

		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&id) 

		if err != nil {
			http.Error(w, "Valid course not entered", http.StatusBadRequest)
			return
		}

		course, err := models.GetCourse(db, id)

		if err != nil {
			http.Error(w, "Get course failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(course)
	}
}

func GetCoursesHandler(db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		if r.Method != http.MethodGet{
			http.Error(w, "Not valid method", http.StatusMethodNotAllowed)
			return
		}

		courses, err := models.GetCourses(db)
		
		if err != nil {
			http.Error(w, "Get courses failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(courses)
	}
}