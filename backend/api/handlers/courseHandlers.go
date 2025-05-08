package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/nynniaw12/ieee-planner/db"
	"github.com/nynniaw12/ieee-planner/scraper"
)

func GetAvailableQuartersHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
		// query distinct quarters
        query := `SELECT DISTINCT quarter FROM courses ORDER BY quarter DESC`

        rows, err := db.Query(query)
        if err != nil {
            http.Error(w, fmt.Sprintf("Error querying quarters: %v", err), http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        var quarters []int
        for rows.Next() {
            var quarter int
            if err := rows.Scan(&quarter); err != nil {
                http.Error(w, fmt.Sprintf("Error scanning quarter: %v", err), http.StatusInternalServerError)
                return
            }
            quarters = append(quarters, quarter)
        }

        if err := rows.Err(); err != nil {
            http.Error(w, fmt.Sprintf("Error iterating rows: %v", err), http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(quarters); err != nil {
            http.Error(w, fmt.Sprintf("Error encoding quarters: %v", err), http.StatusInternalServerError)
            return
        }
    }
}

func GetCoursesByQuarterHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Get quarter parameter from query string
        quarterStr := r.URL.Query().Get("quarter")
        if quarterStr == "" {
            http.Error(w, "Quarter parameter is required", http.StatusBadRequest)
            return
        }

        quarter, err := strconv.Atoi(quarterStr)
        if err != nil {
            http.Error(w, "Invalid quarter format", http.StatusBadRequest)
            return
        }

        // Query for courses by quarter
        query := `SELECT id, title, course_number, topic, overview, url, section, subject, school, quarter 
                 FROM courses WHERE quarter = $1`

        rows, err := db.Query(query, quarter)
        if err != nil {
            http.Error(w, fmt.Sprintf("Error querying courses: %v", err), http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        // Collect courses into a slice
        var courses []*scraper.Course
        courseIDs := make(map[int]*scraper.Course)

        for rows.Next() {
            course := &scraper.Course{}
            var id int
            
            if err := rows.Scan(&id, &course.Title, &course.Number, &course.Topic, 
                              &course.Overview, &course.URL, &course.Section, 
                              &course.Subject, &course.School, &course.Quarter); err != nil {
                http.Error(w, fmt.Sprintf("Error scanning course: %v", err), http.StatusInternalServerError)
                return
            }
            
            courses = append(courses, course)
            courseIDs[id] = course
        }

        if err := rows.Err(); err != nil {
            http.Error(w, fmt.Sprintf("Error iterating rows: %v", err), http.StatusInternalServerError)
            return
        }

        // Load instructors for each course, and meeting times for each course
        for id, course := range courseIDs {
            instructors, err := getInstructorsForCourse(db, id)
            if err != nil {
                http.Error(w, fmt.Sprintf("Error getting instructors: %v", err), http.StatusInternalServerError)
                return
            }
            course.Instructors = instructors

            meetingTimes, err := getMeetingTimesForCourse(db, id)
            if err != nil {
                http.Error(w, fmt.Sprintf("Error getting meeting times: %v", err), http.StatusInternalServerError)
                return
            }
            course.MeetingTimes = meetingTimes
        }

        // Return courses as JSON
        w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(courses); err != nil {
            http.Error(w, fmt.Sprintf("Error encoding courses: %v", err), http.StatusInternalServerError)
            return
        }
    }
}

// Helper function to get instructors for a course
func getInstructorsForCourse(db *sql.DB, courseID int) ([]scraper.Instructor, error) {
    query := `SELECT name, phone, email, office_hours, address 
             FROM instructors WHERE course_id = $1`
             
    rows, err := db.Query(query, courseID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var instructors []scraper.Instructor
    for rows.Next() {
        var instructor scraper.Instructor
        if err := rows.Scan(&instructor.Name, &instructor.Phone, &instructor.Email, 
                          &instructor.OfficeHours, &instructor.Address); err != nil {
            return nil, err
        }
        instructors = append(instructors, instructor)
    }
    
    return instructors, rows.Err()
}

// Helper function to get meeting times for a course
func getMeetingTimesForCourse(db *sql.DB, courseID int) ([]scraper.MeetingTime, error) {
    query := `SELECT location, days, start_time, end_time, time_range 
             FROM meetingtimes WHERE course_id = $1`
             
    rows, err := db.Query(query, courseID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var meetingTimes []scraper.MeetingTime
    for rows.Next() {
        var meetingTime scraper.MeetingTime
        var daysJSON []byte
        
        if err := rows.Scan(&meetingTime.Location, &daysJSON, &meetingTime.StartTime, 
                          &meetingTime.EndTime, &meetingTime.TimeRange); err != nil {
            return nil, err
        }
        
        // Unmarshal JSON days array
        if err := json.Unmarshal(daysJSON, &meetingTime.Days); err != nil {
            return nil, err
        }
        
        meetingTimes = append(meetingTimes, meetingTime)
    }
    
    return meetingTimes, rows.Err()
}

func GetCoursesBySubjectHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Get subject parameter from query string
        subject := r.URL.Query().Get("subject")
        if subject == "" {
            http.Error(w, "Subject parameter is required", http.StatusBadRequest)
            return
        }

        // Query for courses by subject
        query := `SELECT DISTINCT title, course_number, topic, overview, array_agg(quarter) as quarters
                 FROM courses 
                 WHERE subject = $1
                 GROUP BY title, course_number, topic, overview
                 ORDER BY course_number`

        rows, err := db.Query(query, subject)
        if err != nil {
            http.Error(w, fmt.Sprintf("Error querying courses by subject: %v", err), http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        // Collect courses into a slice
        type CourseBySubject struct {
            Title    string `json:"title"`
            Number   string `json:"number"`
            Topic    string `json:"topic"`
            Overview string `json:"overview"`
            Quarters []int  `json:"quarters"`
        }

        var coursesBySubject []*CourseBySubject
        
        for rows.Next() {
            course := &CourseBySubject{}
            var quartersArray []int
            
            if err := rows.Scan(&course.Title, &course.Number, &course.Topic, 
                              &course.Overview, &quartersArray); err != nil {
                http.Error(w, fmt.Sprintf("Error scanning course: %v", err), http.StatusInternalServerError)
                return
            }
            
            course.Quarters = quartersArray
            coursesBySubject = append(coursesBySubject, course)
        }

        // Check for errors from iterating over rows
        if err := rows.Err(); err != nil {
            http.Error(w, fmt.Sprintf("Error iterating rows: %v", err), http.StatusInternalServerError)
            return
        }

        // Return courses as JSON
        w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(coursesBySubject); err != nil {
            http.Error(w, fmt.Sprintf("Error encoding courses: %v", err), http.StatusInternalServerError)
            return
        }
    }
}

func GetCoursesByKeyHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Get key parameter from query string
        key := r.URL.Query().Get("key")
        if key == "" {
            http.Error(w, "Key parameter is required", http.StatusBadRequest)
            return
        }

        // Parse key to get subject and number parts
        parts := strings.Split(key, " ")
        if len(parts) != 2 {
            http.Error(w, "Invalid key format", http.StatusBadRequest)
            return
        }
        
        subject := parts[0]
        numberPattern := parts[1] + "%"

        // Query for courses by key
        query := `SELECT id, title, course_number, topic, overview, url, section, subject, school, quarter 
                 FROM courses 
                 WHERE subject = $1 AND course_number LIKE $2`

        rows, err := db.Query(query, subject, numberPattern)
        if err != nil {
            http.Error(w, fmt.Sprintf("Error querying courses by key: %v", err), http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        // Process results similar to GetCoursesByQuarterHandler
        var courses []*scraper.Course
        courseIDs := make(map[int]*scraper.Course)

        for rows.Next() {
            course := &scraper.Course{}
            var id int
            
            if err := rows.Scan(&id, &course.Title, &course.Number, &course.Topic, 
                              &course.Overview, &course.URL, &course.Section, 
                              &course.Subject, &course.School, &course.Quarter); err != nil {
                http.Error(w, fmt.Sprintf("Error scanning course: %v", err), http.StatusInternalServerError)
                return
            }
            
            courses = append(courses, course)
            courseIDs[id] = course
        }

        // Load related data
        for id, course := range courseIDs {
            course.Instructors, _ = getInstructorsForCourse(db, id)
            course.MeetingTimes, _ = getMeetingTimesForCourse(db, id)
        }

        // Return courses as JSON
        w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(courses); err != nil {
            http.Error(w, fmt.Sprintf("Error encoding courses: %v", err), http.StatusInternalServerError)
            return
        }
    }
}

// GetMajorRequirementsHandler handles requests for major requirements
func GetMajorRequirementsHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		major := r.URL.Query().Get("major")
		if major == "" {
			http.Error(w, "Major parameter is required", http.StatusBadRequest)
			return
		}

		reqs, err := db.GetMajorReqsFromDatabase(database, major)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				http.Error(w, "Major not found", http.StatusNotFound)
			} else {
				http.Error(w, "Server error", http.StatusInternalServerError)
				log.Printf("Error retrieving major requirements: %v", err)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(reqs)
	}
}

// GetAvailableMajorsHandler returns a list of all available majors
func GetAvailableMajorsHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		majors, err := db.GetAllMajorsFromDatabase(database)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			log.Printf("Error retrieving majors: %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(majors)
	}
}