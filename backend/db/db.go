package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

func ConnectToDB() *sql.DB {
	host     := os.Getenv("DB_HOST")
	portstr  := os.Getenv("DB_PORT")
	user     := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname   := os.Getenv("DB_NAME")

	port, err := strconv.Atoi(portstr)
	if err != nil {
		log.Fatal(err)
	}

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}

// CreateCoursesTable drops the table if it exists and creates it fresh
func CreateCoursesTable(db *sql.DB) error {
	// Drop if exists (optional safety)
	dropQuery := `
	DROP TABLE IF EXISTS courses;
	`

	_, err := db.Exec(dropQuery)
	if err != nil {
		return fmt.Errorf("failed to drop table: %w", err)
	}

	// Create new table
	createQuery := `
	CREATE TABLE courses (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		department TEXT NOT NULL,
		professor TEXT NOT NULL,
		time TIME NOT NULL,
		days TEXT NOT NULL,
		location TEXT NOT NULL,
		description TEXT NOT NULL,
		prerequisites JSON NOT NULL
	);
	`

	_, err = db.Exec(createQuery)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}
