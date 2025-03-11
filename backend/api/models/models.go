package models

type Course struct {
	Department string `json:"department"`
	ID	int `json:"id"` // `: how filed is represented in JSON object
	Name string `json:"name"`
	Description string `json:"description"`
	Prerequisites []string `json:"prerequisites"`
}