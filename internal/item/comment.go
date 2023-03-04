package item

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Comment is based off, but does not follow, the structure of the official API.
type Comment struct {
	*Item
	Story      string   `json:"parent"` // Title of the story.
	Text       []string `json:"text"`
	Context    string   `json:"context"`    // Link to the comment in the original thread.
	Discussion string   `json:"discussion"` // Link to the original thread.
}

func (c Comment) Itemize() *Item {
	c.Item.Entry = c.Preview(75)
	return c.Item
}

func (c Comment) Details() {
	fmt.Println("On:", c.Story)
	fmt.Println("Preview:", c.Preview(75))
	fmt.Println("Context:", c.Context)
}

// Preview returns the first `n` characters of the comment text, truncating the last line of text if needed.
// The text will be returned as a single string, separated by spaces.
func (c Comment) Preview(n int) string {
	var text []string
	var totalLength int

	for _, line := range c.Text {
		length := utf8.RuneCountInString(line)
		if length <= n-totalLength {
			// If the line fits completely within the max length (150), then add it directly.
			text = append(text, line)
			totalLength += length
		} else {
			// Otherwise, add only those characters that bring the total count to 150.
			text = append(text, string([]rune(line)[:n-totalLength]))
			break
		}
	}

	return strings.ReplaceAll(strings.Join(text, " "), "\n", "")
}
