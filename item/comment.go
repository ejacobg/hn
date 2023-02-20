package item

// Comment is based off, but does not follow, the structure of the official API.
type Comment struct {
	ID         string   `json:"id"`
	Story      string   `json:"parent"` // Title of the story.
	Text       []string `json:"text"`
	Context    string   `json:"context"`    // Link to the comment in the original thread.
	Discussion string   `json:"discussion"` // Link to the original thread.
	Tags       []Tag    `json:"tags,omitempty"`
	Statuses   []Status `json:"statuses,omitempty"`
}
