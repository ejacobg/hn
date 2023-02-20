package item

import "fmt"

// Story is based off, but does not follow, the structure of the official API.
type Story struct {
	*Item
	Title      string `json:"title"`
	URL        string `json:"url"`
	Discussion string `json:"discussion"`
}

func (s Story) Itemize() *Item {
	s.Item.Entry = s.Title
	return s.Item
}

func (s Story) Details() {
	fmt.Println("Title:", s.Title)
	fmt.Println("URL:", s.URL)
	fmt.Println("Discussion:", s.Discussion)
}
