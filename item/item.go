package item

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

// Itemizers will convert a slice of itemizer values into a slice of itemizer interfaces.
func Itemizers[I Itemizer](items []I) []Itemizer {
	if items == nil {
		return nil
	}
	izs := make([]Itemizer, len(items))
	for i, itm := range items {
		izs[i] = itm
	}
	return izs
}

type Detailer interface {
	// Details should print a brief summary to the screen about the content of an Item.
	Details()
}
