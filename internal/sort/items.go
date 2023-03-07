package sort

import (
	"fmt"
	"github.com/ejacobg/hn/internal/item"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// All valid lines have the form: [<status>] [<entry>](<id>)
var validLine = regexp.MustCompile(`\[(.*)] \[(.*)]\((.*)\)`)

func NewItemIndex[I item.Itemizer](pages []item.Page[I]) *map[string]*item.Item {
	index := make(map[string]*item.Item)

	for i := 0; i < len(pages); i++ {
		for j := range pages[i].Items {
			itm := pages[i].Items[j].Itemize()
			index[itm.ID] = itm
		}
	}

	return &index
}

func categorizeItems(files []string, index *map[string]*item.Item) {
	seen := make(map[string]bool)
	for _, fileName := range files {
		category := strings.TrimSuffix(filepath.Base(fileName), ".md") // Assumes markdown files.

		processLines(fileName, func(line string) {
			matches := validLine.FindStringSubmatch(line)
			if len(matches) > 0 {
				state, _, id := matches[1], matches[2], matches[3] // Entry capture group is not used.
				if itm, ok := (*index)[id]; ok {
					// Apply the current category only upon user request or if the item hasn't been updated before.
					// Otherwise, the item will be left alone.
					// Duplicate entries with the same category are ignored.
					if seen[id] && resolveCategoryConflict(itm.Entry, itm.Category, category) || !seen[id] {
						itm.Category = category
						itm.State = state
					}
					seen[id] = true
				}
			}
		})
	}
}

func resolveCategoryConflict(entry, curr, new string) (keep bool) {
	// If there is no conflict, then return.
	if curr == new {
		return true
	}
	var input rune
	fmt.Printf("%s: category %q conflicts with previous category %q\n", entry, new, curr)
	fmt.Print("[K]eep or [I]gnore: ")
	fmt.Scanf("%c\n", &input)
	for input != 'K' && input != 'I' {
		fmt.Print("[K]eep or [I]gnore: ")
		fmt.Scanf("%c\n", &input)
	}
	return input == 'K'
}

// removeDuplicateItems will read through all the given category files and remove all duplicate or incorrect entries.
// Note: the sorting operation may create new category files. Make sure to pass in all category files to this function or else they will not be cleaned up.
func removeDuplicateItems(categories []string, index *map[string]*item.Item) {
	for _, fileName := range categories {
		// keep holds lines that will be preserved.
		var keep []string

		category := strings.TrimSuffix(filepath.Base(fileName), ".md") // Assumes markdown files.

		processLines(fileName, func(line string) {
			matches := validLine.FindStringSubmatch(line)
			if len(matches) > 0 {
				id := matches[3]
				// If the item belongs in this category, keep its entry.
				if itm, ok := (*index)[id]; ok && itm.Category == category {
					keep = append(keep, line)
					// By removing this item from the index, we prevent multiple lines from the same file being kept.
					// Performing this deletion also makes adding missing entries easier.
					delete(*index, id)
				}
			} else {
				// If the regex doesn't match, keep the line anyway since it might be used for formatting.
				keep = append(keep, line)
			}
		})

		// Open the file again, this time truncating it.
		file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			fmt.Println("Error opening file for writing:", err)
			continue
		}

		// Write back all the lines that we have decided to keep.
		for _, line := range keep {
			_, err := file.WriteString(line + "\n")
			if err != nil {
				fmt.Println("Error writing to file:", err)
			}
		}

		file.Close()
	}
}

func addMissingItems(directory string, index *map[string]*item.Item) {
	for _, itm := range *index {
		if itm.Category == "" {
			continue
		}

		// Append to category file.
		path := directory + "/" + itm.Category + ".md"
		entry := fmt.Sprintf("[%s] [%s](%s)\n", itm.State, itm.Entry, itm.ID)
		err := writeString(path, os.O_APPEND|os.O_CREATE, entry)
		if err != nil {
			fmt.Println("Error writing to file:", err)
		}
	}
}
