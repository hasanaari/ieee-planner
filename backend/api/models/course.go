package models

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type Course struct {
	ID int
	Name string
	Department string
	Professor string
	Time time.Time
	Days string
	Location string
}

func AddCourse(db *sql.DB, course Course) error {
	
	query := "INSERT INTO courses (id, name, department, professor, time, days, location) VALUES ($1, $2, $3, $4, $5, $6, $7)"

	_, err := db.Exec(query,
		course.ID,
		course.Name,
		course.Department,
		course.Professor,
		course.Time,
		course.Days,
		course.Location)	

	if err != nil {
		return err
	}

	return err
}

func EditCourse (db *sql.DB, course Course) error {
	query := "UPDATE courses SET name = $1, department = $2, professor = $3, time = $4, days = $5, location = $6 WHERE id = $7"

	_, err := db.Exec(query, 
		course.Name,
		course.Department,
		course.Professor,
		course.Time,
		course.Days,
		course.Location,
		course.ID)

	if err != nil {
		return err
	}

	return err
}

func RemoveCourse (db *sql.DB, id int) error {
	query := "DELETE FROM courses WHERE id=$1"

	_, err := db.Exec(query,
	id)

	if err != nil {
		return err
	}
	
	return err
}

func ClearCourses (db *sql.DB) error {
	query := "DELETE FROM courses"

	_, err := db.Exec(query)

	if err != nil {
		return err
	}

	return err
}

func GetCourse (db *sql.DB, id int) (Course, error){
	query := "SELECT * FROM courses WHERE id=$1"

	row := db.QueryRow(query,
	id)
	
	var selectedCourse Course

	err := row.Scan(&selectedCourse.ID, 
		&selectedCourse.Name, 
		&selectedCourse.Department, 
		&selectedCourse.Professor, 
		&selectedCourse.Time, 
		&selectedCourse.Days, 
		&selectedCourse.Location)

	if err != nil {
		return Course{}, err
	}

	return selectedCourse, err
}

func GetCourses (db *sql.DB) ([]Course, error) {
	query := "SELECT * from COURSES"

	rows, err := db.Query(query)

	courses := make([]Course, 0)

	if err != nil {
		return nil, err
	}

	for rows.Next(){
		var selectedCourse Course
		err := rows.Scan(&selectedCourse.ID, 
			&selectedCourse.Name, 
			&selectedCourse.Department, 
			&selectedCourse.Professor, 
			&selectedCourse.Time, 
			&selectedCourse.Days, 
			&selectedCourse.Location)
		
		if err != nil {
			log.Fatal(err)
		}

		courses = append(courses, selectedCourse)
	}

	return courses, nil
}