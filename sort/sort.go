package sort

import (
	"fmt"
	"github.com/ejacobg/hn/item"
	"os"
)

type Sorter interface {
	item.Detailer
	item.Itemizer
}

func Items[I Sorter](directory string) error {
	// Read exports and category files from correct directory.
	categories, pages, err := item.ReadDirectory[I](directory)
	if err != nil {
		return fmt.Errorf("items: failed to read directory: %w", err)
	}

	// Read in reading statuses and categories.
	index := NewItemIndex(pages)
	categorizeItems(categories, index)

	// Show CLI
	var input rune
Loop:
	for i := 0; i < len(pages); i++ {
		for j := 0; j < len(pages[i].Items); j++ {
			sorter := pages[i].Items[j]
			itm := sorter.Itemize()

			// If the item already has a category, skip it.
			if itm.Category != "" {
				continue
			}

			// Otherwise, let user determine the category.
			sorter.Details()
			fmt.Println("1 - Advice")
			fmt.Println("2 - Career")
			fmt.Println("3 - Hobby")
			fmt.Println("4 - Money")
			fmt.Println("5 - Opinion")
			fmt.Println("6 - Product")
			fmt.Println("7 - Project")
			fmt.Println("8 - Reference")
			fmt.Println("9 - Skill")
			fmt.Println("0 - Tip")
			fmt.Println("f - Finish")
			fmt.Scanf("%c\n", &input)
			switch input {
			case '1':
				itm.Category = item.Advice
			case '2':
				itm.Category = item.Career
			case '3':
				itm.Category = item.Hobby
			case '4':
				itm.Category = item.Money
			case '5':
				itm.Category = item.Opinion
			case '6':
				itm.Category = item.Product
			case '7':
				itm.Category = item.Project
			case '8':
				itm.Category = item.Reference
			case '9':
				itm.Category = item.Skill
			case '0':
				itm.Category = item.Tip
			case 'f':
				break Loop
			default:
				continue
			}

			// Append new category to file.
			path := directory + "/" + itm.Category + ".md"
			entry := fmt.Sprintf("[%s] [%s](%s)\n", itm.State, itm.Entry, itm.ID)
			err := writeString(path, os.O_APPEND|os.O_CREATE, entry)
			if err != nil {
				fmt.Println("Error writing to file:", err)
				continue
			}
		}
	}

	// Write updated data back into export files.
	for i := range pages {
		err := pages[i].Write()
		if err != nil {
			fmt.Println("Error writing to file:", err)
		}
	}

	// Perform cleanup.
	removeDuplicateItems(categories, index)
	addMissingItems(directory, index)

	return nil
}
