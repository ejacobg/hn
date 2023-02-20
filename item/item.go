package item

import (
	"encoding/json"
	"os"
)

// Item represents the smallest unit of organization within the system.
// In order for a type to be sortable, it must embed an Item.
type Item struct {
	ID       string `json:"id"`
	Category string `json:"category,omitempty"`
	State    string `json:"state,omitempty"`
	Entry    string `json:"-"`
}

// Story.Itemize() and Comment.Itemize() both need to have value receivers or else a []Page[Story|Item] cannot be instantiated.
// In order to share the underlying Item while still using a value receiver, a *Item needs to be embedded.

// Alternatively, you may instantiate a []Page[*Story|*Item], which would allow embedding of Item (instead of *Item).
// In this case, you would have to use a pointer receiver for the given methods.

type Itemizer interface {
	// Itemize should return the object's embedded Item in order to be edited by the sorter.
	Itemize() *Item
}

type Detailer interface {
	// Details should print a brief summary to the screen about the content of an Item.
	Details()
}

// FromList parses a given JSON file and returns a slice of its contents.
// The content of the given JSON file should be an array of objects (e.g. an array of Story).
func FromList[I Itemizer](path string) ([]I, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var i []I
	err = json.Unmarshal(file, &i)
	if err != nil {
		return nil, err
	}
	return i, nil
}
