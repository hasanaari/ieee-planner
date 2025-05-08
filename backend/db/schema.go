// Creates tables if they do not already exist
package db

import (
	"database/sql"
	"fmt"
)

// CreateMajorReqsTable creates the table for storing major requirements
func CreateMajorReqsTableIfNotExists(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS major_requirements (
		id SERIAL PRIMARY KEY,
		major TEXT UNIQUE NOT NULL,
		is_engineering BOOLEAN NOT NULL,
		requirements JSONB NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating major requirements table: %w", err)
	}
	
	fmt.Println("Major requirements table created or already exists")
	return nil
}


// Creates courses table if it doesn't exist
func CreateCoursesTableIfNotExists (db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS courses (
	id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    course_number VARCHAR(50) NOT NULL,
    topic TEXT, 
    overview TEXT, 
    url VARCHAR(255), 
    section INTEGER,
    subject VARCHAR(100),
    school VARCHAR(100), 
    quarter INTEGER NOT NULL
	)`

	_, err := db.Exec(query)

	if err != nil {
		return fmt.Errorf("failed to create courses table: %w", err)
	}

	return nil
}

// Creates instructors table if it doesn't exist
func CreateInstructorsTableIfNotExists (db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS Instructors (
    id SERIAL PRIMARY KEY, 
    course_id INTEGER NOT NULL REFERENCES Courses(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    email VARCHAR(100),
    office_hours TEXT,
    address TEXT
	)`

	_, err := db.Exec(query)

	if err != nil {
		return fmt.Errorf("failed to create instructors table: %w", err)
	}

	return nil
}

// Creates meetingtimes table if it doesn't exist
func CreateMeetingTimesTableIfNotExists (db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS MeetingTimes (
    id SERIAL PRIMARY KEY,
    course_id INTEGER NOT NULL REFERENCES Courses(id) ON DELETE CASCADE,
    location VARCHAR(255) NOT NULL,
    days JSONB NOT NULL,
    start_time TIMESTAMP,
    end_time TIMESTAMP,
    time_range VARCHAR(100) NOT NULL
	)`

	_, err := db.Exec(query)

	if err != nil {
		return fmt.Errorf("failed to create meeting times table: %w", err)
	}

	return nil
}

// Creates all necessary database tables for MVP
func InitializeTables (db *sql.DB) error {
	err := CreateCoursesTableIfNotExists(db)

	if err != nil {
		return err
	}

	err = CreateInstructorsTableIfNotExists(db)

	if err != nil {
		return err
	}

	err = CreateMeetingTimesTableIfNotExists(db)

	if err != nil {
		return err
	}

	err = CreateMajorReqsTableIfNotExists(db)

	if err != nil {
		return err
	}

	return nil
}