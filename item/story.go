package item

// Story is based off, but does not follow, the structure of the official API.
type Story struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	URL        string   `json:"url"`
	Discussion string   `json:"discussion"`
	Tags       []Tag    `json:"tags,omitempty"`
	Statuses   []Status `json:"statuses,omitempty"`
}
