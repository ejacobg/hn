package sort

import (
	"fmt"
	"github.com/ejacobg/hn/item"
	"path/filepath"
)

// readDirectory will parse the export and category files of the given directory, and return a list of category files, and a parsed set of pages.
func readDirectory[I item.Itemizer](directory string) ([]string, []item.Page[I], error) {
	exports, err := filepath.Glob(directory + "/exported/*.json")
	if err != nil {
		return nil, nil, fmt.Errorf("readDirectory: failed to find export files: %w", err)
	}

	categories, err := filepath.Glob(directory + "/*.md")
	if err != nil {
		return nil, nil, fmt.Errorf("readDirectory: failed to find category files: %w", err)
	}

	pages := make([]item.Page[I], len(exports))
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
