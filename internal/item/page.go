package item

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

// ReadDirectory will parse the export and category files of the given directory, and return a list of category files, and a parsed set of pages.
func ReadDirectory[I Itemizer](directory string) ([]string, []Page[I], error) {
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
		pages[i].Items, err = FromList[I](file)
		if err != nil {
			return nil, nil, fmt.Errorf("readDirectory: failed to unmarshal file %s: %w", file, err)
		}
	}

	return categories, pages, nil
}
