package models

type Course struct {
	Department string `json:"department"`
	ID	int `json:"id"` // `: how filed is represented in JSON object
	Name string `json:"name"`
	Description string `json:"description"`
	Prerequisites []string `json:"prerequisites"`
}

// ClickAction represents a click action with the element ID and other details.
type ClickAction struct {
    Action  string `json:"action"`
    Element int    `json:"element"`
}

// TypeAction represents a typing action with the element ID and text to type.
type TypeAction struct {
    Action  string `json:"action"`
    Element int    `json:"element"`
    Text    string `json:"text"`
}

// RememberInfoFromSite represents an action to remember information from the site.
type RememberInfoFromSite struct {
    Action string `json:"action"`
    Info   string `json:"info"`
}

// Done represents the action that signifies the task is completed.
type Done struct {
    Action string `json:"action"`
}