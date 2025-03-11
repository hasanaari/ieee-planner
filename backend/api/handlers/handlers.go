package handlers

import (
	"encoding/json" // ecnode/deconde, sending/recieving JSON data
	"fmt"
	"net/http" // handling HTTP requests and responses
	"ieee-planner/backend/api/models"
)

// w: write response back to client
// r: incoming HTTP request
// http is from net/http package
// func FunctionName(param1 Type1, param2 Type2)
func CoursesHandler(w http.ResponseWriter, r *http.Request){
	// initialise courses (array of Course structs)
	courses := []Course{ //:= is variable declaration
		{Department: "CS", ID: 111, Name:"Fundamentals of Computer Programming I", Description: "Fundamentals of CS", Prerequisites: []string{}},
		{Department: "CS", ID: 150, Name:"Fundamentals of Computer Programming 1.5", Description: "Fundamentals of CS 2", Prerequisites: []string{"Fundamentals of Computer Programming I"}},
	}
	// Header: additional context about request/response
	w.Header().Set("Content-Type", "application/json") // sets header of response to JSON format

	encoder := json.NewEncoder(w) // JSON-encoded data, writes to w
	encoder.Encode(courses) // convert courses into JSON, writes to w
}

func SchedulesHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "This would return a schedule planner")
}