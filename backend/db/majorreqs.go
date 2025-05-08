package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/nynniaw12/ieee-planner/scraper"
)

// WriteMajorReqsToDatabase writes a single major's requirements to the database
func WriteMajorReqsToDatabase(db *sql.DB, mr *scraper.MajorRequirements) error {
	// Convert the requirements to JSON
	reqsJSON, err := json.Marshal(mr.AllRequirements)
	if err != nil {
		return fmt.Errorf("error converting requirements to JSON: %w", err)
	}

	// Upsert query to either insert or update existing major
	query := `
	INSERT INTO major_requirements (major, is_engineering, requirements)
	VALUES ($1, $2, $3)
	ON CONFLICT (major) 
	DO UPDATE SET 
		is_engineering = $2,
		requirements = $3,
		updated_at = CURRENT_TIMESTAMP
	RETURNING id;`

	var id int
	err = db.QueryRow(query, strings.ToLower(mr.Major), mr.IsEngineering, reqsJSON).Scan(&id)
	if err != nil {
		return fmt.Errorf("error inserting/updating major requirements: %w", err)
	}

	log.Printf("Major %s requirements saved with ID %d\n", mr.Major, id)
	return nil
}

// WriteBulkMajorReqsToDatabase writes multiple majors' requirements to the database
func WriteBulkMajorReqsToDatabase(db *sql.DB, majors map[string]*scraper.MajorRequirements) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback() // Will be ignored if transaction is committed

	// Prepare the statement for better performance with multiple inserts
	stmt, err := tx.Prepare(`
		INSERT INTO major_requirements (major, is_engineering, requirements)
		VALUES ($1, $2, $3)
		ON CONFLICT (major) 
		DO UPDATE SET 
			is_engineering = $2,
			requirements = $3,
			updated_at = CURRENT_TIMESTAMP
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	for _, mr := range majors {
		reqsJSON, err := json.Marshal(mr.AllRequirements)
		if err != nil {
			return fmt.Errorf("error converting requirements to JSON for major %s: %w", mr.Major, err)
		}

		_, err = stmt.Exec(strings.ToLower(mr.Major), mr.IsEngineering, reqsJSON)
		if err != nil {
			return fmt.Errorf("error inserting/updating major %s: %w", mr.Major, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	log.Printf("Successfully saved %d major requirements to database\n", len(majors))
	return nil
}

// GetMajorReqsFromDatabase retrieves a single major's requirements from the database
func GetMajorReqsFromDatabase(db *sql.DB, major string) (*scraper.MajorRequirements, error) {
	query := `
	SELECT major, is_engineering, requirements
	FROM major_requirements
	WHERE major = $1;`

	var majorName string
	var isEngineering bool
	var reqsJSON []byte

	err := db.QueryRow(query, strings.ToLower(major)).Scan(&majorName, &isEngineering, &reqsJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("major %s not found", major)
		}
		return nil, fmt.Errorf("error retrieving major requirements: %w", err)
	}

	var allReqs []any
	err = json.Unmarshal(reqsJSON, &allReqs)
	if err != nil {
		return nil, fmt.Errorf("error parsing requirements JSON: %w", err)
	}

	// Construct the MajorRequirements object
	majorReqs := &scraper.MajorRequirements{
		Major:           majorName,
		IsEngineering:   isEngineering,
		AllRequirements: allReqs,
	}

	return majorReqs, nil
}

// GetAllMajorsFromDatabase retrieves all majors from the database
func GetAllMajorsFromDatabase(db *sql.DB) ([]string, error) {
	query := `
	SELECT major
	FROM major_requirements
	ORDER BY major;`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying majors: %w", err)
	}
	defer rows.Close()

	var majors []string
	for rows.Next() {
		var major string
		if err := rows.Scan(&major); err != nil {
			return nil, fmt.Errorf("error scanning major: %w", err)
		}
		majors = append(majors, major)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating majors: %w", err)
	}

	return majors, nil
}

// InitializeMajorReqsFromScraper loads major requirements from the scraper and stores them in the database
func InitializeMajorReqsFromScraper(db *sql.DB, path string) error {
	// Load major requirements from files
	store, err := scraper.NewMajorRequirementsStore(path)
	if err != nil {
		return fmt.Errorf("error creating major requirements store: %w", err)
	}

	// Write all requirements to database
	return WriteBulkMajorReqsToDatabase(db, store.RequirementsByMajor)
}