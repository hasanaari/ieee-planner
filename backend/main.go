package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
	"ieee-planner/backend/api/handlers"
)

// Based on IEEE #2
type Response struct{
    BriefExplanation string `json: "briefExplanation"`
    NextAction  struct{
        Action string `json:"action"`
        Element int `json:"element"`
    }   `json:"nextAction"`
}

type OpenAIRequest struct{
    Model string `json:"model"`
    Prompt string `json:"prompt"`
    MaxTokens int `json:"max_tokens"`
    Temperature float64 `json:"temperature"`
}

func main(){

	// testing handlers
	http.HandleFunc("/api/courses", handlers.CoursesHandler) // request to /courses, call CoursesHandler
    http.HandleFunc("/api/schedules", handlers.SchedulesHandler)
    fmt.Println("Server starting on :8080...")
    log.Fatal(http.ListenAndServe(":8080", nil)) // error will stop program

    prompt := `task: Find out how to get course information for CS101`
    
    // previous action
    // click on the course catalog
    
    // instructions
    // observe the webpage, and suggest the next action

    // Define the OpenAI API request
	requestBody := OpenAIRequest{
		Model:       "text-davinci-003", // You can use other models like gpt-3.5-turbo or gpt-4 as well
		Prompt:      prompt,
		MaxTokens:   150,
		Temperature: 0.7,
	}

	// Send the request to OpenAI API
	resp, err := sendRequestToOpenAI(requestBody)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}

	// Print the response
	fmt.Printf("Response from OpenAI: %v\n", resp)
}

func sendRequestToOpenAI(request OpenAIRequest) (*Response, error) {
	// Convert request body to JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request data: %v", err)
	}

	// Send POST request to OpenAI API
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	// Set the authorization header with the API key
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openAIAPIKey)

	// Perform the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Read and parse the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// If the status is not 200, handle the error
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response from OpenAI: %s", body)
	}

	// Parse the response into the struct
	var openAIResponse Response
	err = json.Unmarshal(body, &openAIResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response body: %v", err)
	}

	return &openAIResponse, nil
}
