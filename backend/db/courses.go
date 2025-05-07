package db

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/nynniaw12/ieee-planner/scraper"
)

// Writes course data -> courses table, instructor data -> instructors table, meetingtimes -> meetingtimes table
func WriteCourseDataToDatabase (db *sql.DB, course scraper.Course) error {
	tx, err := db.Begin()

	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer tx.Rollback()

	query := `INSERT INTO courses (title, course_number, topic, overview, url, section, subject, school, quarter)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			  RETURNING id`

	var courseID int 

	err = tx.QueryRow(query, course.Title,
	course.Number,
	course.Topic,
	course.Overview,
	course.URL, 
	course.Section, 
	course.Subject, 
	course.School, 
	course.Quarter).Scan(&courseID)

	if err != nil {
		return fmt.Errorf("failed to write course to database: %w", err)
	}

	query = `INSERT INTO instructors (course_id, name, phone, email, office_hours, address)
	VALUES ($1, $2, $3, $4, $5, $6)`

	for _, instructor := range course.Instructors {
		_, err = tx.Exec(query, courseID, instructor.Name, instructor.Phone, instructor.Email, instructor.OfficeHours, instructor.Address)

		if err != nil {
			return fmt.Errorf("failed to write instructor to database: %w", err)
		}
	}

	query = `INSERT INTO meetingtimes (course_id, location, days, start_time, end_time, time_range)
	VALUES ($1, $2, $3, $4, $5, $6)`

	for _, meetingtime := range course.MeetingTimes {

		daysJSON, err := json.Marshal(meetingtime.Days)

    	if err != nil {
        return fmt.Errorf("failed to marshal days to JSON: %w", err)
    	}

		_, err = tx.Exec(query, courseID, meetingtime.Location, daysJSON, meetingtime.StartTime, meetingtime.EndTime, meetingtime.TimeRange)

		if err != nil {
			return fmt.Errorf("failed to write meetingtime to database: %w", err)
		}
	}

	err = tx.Commit()

	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}