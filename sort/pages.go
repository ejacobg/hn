package sort

import (
	"encoding/json"
	"fmt"
	"github.com/ejacobg/hn/item"
	"os"
	"path/filepath"
)

type Page[I item.Itemizer] struct {
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

// readDirectory will parse the export and category files of the given items, and return a list of category files, and a parsed set of pages.
func readDirectory[I item.Itemizer](saveType, itemType string) ([]string, []Page[I], error) {
	directory := "./" + saveType + "/" + itemType
	exports, err := filepath.Glob(directory + "/exported/*.json")
	if err != nil {
		return nil, nil, fmt.Errorf("readDirectory: failed to find export files: %w", err)
	}

	categories, err := filepath.Glob(directory + "/*.md")
	if err != nil {
		return nil, nil, fmt.Errorf("readDirectory: failed to find category files: %w", err)
	}

	pages := make([]Page[I], len(exports))
	for i, file := range exports {
		var err error
		pages[i].Path = file
		pages[i].Items, err = item.FromList[I](file)
		if err != nil {
			return nil, nil, fmt.Errorf("readDirectory: failed to unmarshal file %s: %w", file, err)
		}
	}

	return categories, pages, nil
}
