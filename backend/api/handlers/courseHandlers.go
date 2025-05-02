package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/nynniaw12/ieee-planner/api/models"
)

// Add course
func AddCourseHandler(db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		id := r.PathValue("id")
		
		decoder := json.NewDecoder(r.Body)

		var course models.Course

		err := decoder.Decode(&course)

		if err != nil {
			http.Error(w, "Valid course not entered", http.StatusBadRequest)
			return
		}

		courseID, err := strconv.Atoi(id)
		
		if err != nil {
			http.Error(w, "Invalid ID entered", http.StatusBadRequest)
			return
		}

		course.ID = courseID

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

		id := r.PathValue("id")

		decoder := json.NewDecoder(r.Body)

		var course models.Course

		err := decoder.Decode(&course)

		if err != nil {
			http.Error(w, "Valid course not entered", http.StatusBadRequest)
			return
		}

		courseID, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid ID entered", http.StatusBadRequest)
			return
		}

		course.ID = courseID 

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

		id := r.PathValue("id")

		courseID, err := strconv.Atoi(id)

		if err != nil {
			http.Error(w, "Invalid ID entered", http.StatusBadRequest)
			return
		}

		err = models.RemoveCourse(db, courseID)

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

		id := r.PathValue("id")

		courseID, err := strconv.Atoi(id)

		if err != nil {
			http.Error(w, "Invalid ID entered", http.StatusBadRequest)
			return
		}

		course, err := models.GetCourse(db, courseID)

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