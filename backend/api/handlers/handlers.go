package handlers

import (
	"encoding/json" // ecnode/deconde, sending/recieving JSON data
	"fmt"
	"net/http" // handling HTTP requests and responses
    "io/ioutil"
    "os"
    "bytes"
	
	"github.com/nynniaw12/ieee-planner/backend/api/models"
)

// Prompt template
var prompt_template = `task: Find out how to parse course listings to see what courses are available.

type ClickAction = %v
type TypeAction = %v
type RememberInfoFromSite = %v
type Done = %v

## response format
{
  "briefExplanation": string,
  "nextAction": ClickAction | TypeAction | RememberInfoFromSite | Done
}

## response examples
{
  "briefExplanation": "I'll click on the course catalog to view available courses.",
  "nextAction": { "action": "click", "element": 5 }
}
{
  "briefExplanation": "I'll type 'Northwestern University CS101' into the search bar to find specific course details.",
  "nextAction": { "action": "type", "element": 8, "text": "Northwestern University CS101" }
}
{
  "briefExplanation": "The course description mentions prerequisites. I'll remember this information for later reference.",
  "nextAction": { "action": "remember-info", "info": "CS101 requires MATH 220 as a prerequisite." }
}
{
  "briefExplanation": "I have extracted course details, including prerequisites and times.",
  "nextAction": { "action": "done" }
}

## previous action
%s

## stored info
%s

## instructions
# observe the screenshot, and think about the next action
# prioritize intuitive actions i.e. clicking on elements located above on the screen
# do not propose actions based on speculation
# done when user input is necessary
# output your response in a json markdown code block
`
var courseListingURL = "https://catalogs.northwestern.edu/undergraduate/courses-az/"

// w: write response back to client
// r: incoming HTTP request
// http is from net/http package
// func FunctionName(param1 Type1, param2 Type2)
func CoursesHandler(w http.ResponseWriter, r *http.Request){
	// initialise courses (array of Course structs)
	courses := []models.Course{ //:= is variable declaration
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

func sendRequestToOpenAI(apiKey, prompt string) (string, error) {
    url := "https://api.openai.com/v1/completions"

    requestBody, err := json.Marshal(map[string]interface{}{
        "model":       "gpt-4",
        "prompt":      prompt,
        "max_tokens":  150,
        "temperature": 0.7,
    })
    if err != nil {
        return "", fmt.Errorf("error marshaling request: %v", err)
    }

	// huh, what are post requests
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
    if err != nil {
        return "", fmt.Errorf("error creating request: %v", err)
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+apiKey)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", fmt.Errorf("error sending request: %v", err)
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("error reading response: %v", err)
    }

    return string(body), nil
}

func GenerateHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    apiKey := os.Getenv("OPENAI_API_KEY")
    if apiKey == "" {
        http.Error(w, "Missing API key", http.StatusInternalServerError)
        return
    }

    // Parse request body for task, prev action, and stored info
    var requestData struct {
        Task string `json:"task"`
        Prev string `json:"prev"`
        Info string `json:"info"`
    }

    err := json.NewDecoder(r.Body).Decode(&requestData)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Format the prompt with user input
    prompt := fmt.Sprintf(prompt_template, ClickAction, TypeAction, RememberInfoFromSite, Done, requestData.Task, requestData.Prev, requestData.Info, courseListingURL)

    // Call OpenAI API
    response, err := sendRequestToOpenAI(apiKey, prompt)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
        return
    }

    // Send response back to client
    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(response))
}
