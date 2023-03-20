package update

import (
	"fmt"
	"github.com/ejacobg/hn/internal/item"
)

// Shuffle will redistribute the items across all the given pages. It assumes that the pages are given in chronological order, with the oldest pages first.
// Each page should have at most 30 items, and pages at the beginning of the list are allowed to take items from pages toward the end of the list.
func Shuffle[I item.Itemizer](pages []item.Page[I], limit int) {
	for i := 0; i < len(pages); i++ {
		length := len(pages[i].Items)
		switch {
		// If there are less than 30 items, take from the other pages until you have 30.
		case length < limit:
			take := limit - length
			for j := i + 1; j < len(pages); j++ {
				// If there aren't enough items to take, then take everything and check the next page.
				if len(pages[j].Items) < take {
					// Copy items over (this operation does not affect pages[j].Items).
					// Items from the next page should be *prepended* to this page.
					pages[i].Items = append(pages[j].Items, pages[i].Items...)

					// Update our counter.
					take -= len(pages[j].Items)

					// Remove the items from the other slice.
					// Using an empty slice so that the JSON will marshal as [].
					pages[j].Items = []I{}
				} else {
					// If there are too many items in the slice, copy only what we need.
					// Copy the last items from the slice, and prepend them to the beginning of this one.
					pages[i].Items = append(pages[j].Items[len(pages[j].Items)-take:], pages[i].Items...)

					// Remove the items that we've taken.
					pages[j].Items = pages[j].Items[:len(pages[j].Items)-take]

					// We've filled the slice at this point, so we can exit the loop.
					break
				}
			}
		// If there are exactly 30 items, then this page is good.
		case length == limit:
			continue
		// If there are more than 30 items, then dump the excess into the next page (if possible).
		case length > limit:
			// If this is the last page, just keep the items here.
			if i+1 == len(pages) {
				continue
			}

			// Otherwise, attach any excess to the beginning of the next page.
			pages[i+1].Items = append(pages[i].Items[limit:], pages[i+1].Items...)

			// Keep the first 30 elements.
			pages[i].Items = pages[i].Items[:limit]
		}
	}
}

func ShuffleDir[I item.Itemizer](directory string, limit int) {
	_, pages, err := item.ReadDirectory[I](directory)
	if err != nil {
		fmt.Printf("ShuffleDir: failed to read directory %q: %v\n", directory, err)
	}

	Shuffle[I](pages, limit)

	for i := range pages {
		err := pages[i].Write()
		if err != nil {
			fmt.Println("Error writing to file:", err)
		}
	}
}
