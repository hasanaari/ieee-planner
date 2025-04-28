package models

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type Course struct {
	ID            int      `json:"id"`
	Name          string   `json:"name"`
	Department    string   `json:"department"`
	Professor     string   `json:"professor"`
	Time          time.Time `json:"time"`
	Days          string   `json:"days"`
	Location      string   `json:"location"`
	Description   string   `json:"description"`
	Prerequisites []string `json:"prerequisites"`
}

func AddCourse(db *sql.DB, course Course) error {
	// We'll store prerequisites as a JSON array in the database
	prereqsJSON, err := json.Marshal(course.Prerequisites)
	if err != nil {
		return err
	}

	query := `
	INSERT INTO courses (id, name, department, professor, time, days, location, description, prerequisites)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err = db.Exec(query,
		course.ID,
		course.Name,
		course.Department,
		course.Professor,
		course.Time,
		course.Days,
		course.Location,
		course.Description,
		string(prereqsJSON), // store as JSON text
	)

	return err
}

func EditCourse(db *sql.DB, course Course) error {
	prereqsJSON, err := json.Marshal(course.Prerequisites)
	if err != nil {
		return err
	}

	query := `
	UPDATE courses
	SET name = $1, department = $2, professor = $3, time = $4, days = $5, location = $6, description = $7, prerequisites = $8
	WHERE id = $9
	`

	_, err = db.Exec(query,
		course.Name,
		course.Department,
		course.Professor,
		course.Time,
		course.Days,
		course.Location,
		course.Description,
		string(prereqsJSON),
		course.ID,
	)

	return err
}

func RemoveCourse(db *sql.DB, id int) error {
	query := "DELETE FROM courses WHERE id=$1"

	_, err := db.Exec(query, id)

	return err
}

func ClearCourses(db *sql.DB) error {
	query := "DELETE FROM courses"

	_, err := db.Exec(query)

	return err
}

func GetCourse(db *sql.DB, id int) (Course, error) {
	query := "SELECT id, name, department, professor, time, days, location, description, prerequisites FROM courses WHERE id=$1"

	row := db.QueryRow(query, id)

	var selectedCourse Course
	var prereqsJSON string

	err := row.Scan(
		&selectedCourse.ID,
		&selectedCourse.Name,
		&selectedCourse.Department,
		&selectedCourse.Professor,
		&selectedCourse.Time,
		&selectedCourse.Days,
		&selectedCourse.Location,
		&selectedCourse.Description,
		&prereqsJSON,
	)

	if err != nil {
		return Course{}, err
	}

	// Unmarshal prerequisites JSON back to slice
	err = json.Unmarshal([]byte(prereqsJSON), &selectedCourse.Prerequisites)
	if err != nil {
		return Course{}, err
	}

	return selectedCourse, nil
}

func GetCourses(db *sql.DB) ([]Course, error) {
	query := "SELECT id, name, department, professor, time, days, location, description, prerequisites FROM courses"

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	courses := make([]Course, 0)

	for rows.Next() {
		var selectedCourse Course
		var prereqsJSON string

		err := rows.Scan(
			&selectedCourse.ID,
			&selectedCourse.Name,
			&selectedCourse.Department,
			&selectedCourse.Professor,
			&selectedCourse.Time,
			&selectedCourse.Days,
			&selectedCourse.Location,
			&selectedCourse.Description,
			&prereqsJSON,
		)
		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal([]byte(prereqsJSON), &selectedCourse.Prerequisites)
		if err != nil {
			log.Fatal(err)
		}

		courses = append(courses, selectedCourse)
	}

	return courses, nil
}