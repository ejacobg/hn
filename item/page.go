package item

import (
	"encoding/json"
	"os"
)

type Page[I Itemizer] struct {
	Path  string
	Items []I
}

// Write will write the contents of the Items slice into the file given by Path.
func (p *Page[I]) Write() error {
	file, err := os.OpenFile(p.Path, os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "\t")
	return enc.Encode(p.Items)
}
