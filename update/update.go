package update

import (
	"fmt"
	"github.com/ejacobg/hn/auth"
	"github.com/ejacobg/hn/export"
	"github.com/ejacobg/hn/item"
	"github.com/ejacobg/hn/sort"
	"strconv"
	"time"
)

func Items[I item.Itemizer](directory, saveType, itemType string, user auth.User) error {
	// Read exports and category files from correct directory.
	_, pages, err := item.ReadDirectory[I](directory)
	if err != nil {
		return fmt.Errorf("submissions: failed to read directory: %w", err)
	}

	// Obtain an index of all the items present in the directory.
	index := sort.NewItemIndex(pages)

	// Transform the pages slice into a generic form.
	// itemizers := make([]item.Page[item.Itemizer], len(pages))
	// for i := 0; i < len(pages); i++ {
	// 	page := item.Page[item.Itemizer]{
	// 		Path:  pages[i].Path,
	// 		Items: item.Itemizers(pages[i].Items),
	// 	}
	// 	itemizers = append(itemizers, page)
	// }

	var client *export.Client
	switch saveType {
	case "favorite":
		// URL route uses "favorites"
		saveType = "favorites"
		client = export.DefaultClient(user)
	// Upvoted posts require an authenticated client.
	case "upvoted":
		client, err = export.NewClient(user)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("submissions: unrecognized save type: %q", saveType)
	}

	// Begin exporting pages.
	var exported []item.Page[item.Itemizer]
	for page, more := 1, true; more; page++ {
		fmt.Println("page", page)
		items, err := client.GetPage(saveType, itemType, page)
		if err != nil {
			return err
		}

		// If no items could be parsed, then we should not continue trying.
		if items == nil {
			break
		}

		fmt.Println(items)

		// Find the first item that we have saved before.
		save := len(items) // If all the items are unseen, save all of them.
		for i := range items {
			id := items[i].Itemize().ID
			if _, seen := (*index)[id]; seen {
				// Record the index of the first seen item.
				save = i

				// All items after this have been seen before, so there is no point in searching through more pages.
				more = false
				break
			}
		}

		// When we spot the first item we've seen before, then everything before it will be unseen.
		// Save all of these unseen items.
		exported = append(exported, item.Page[item.Itemizer]{Items: items[:save]})

		// Add all the new items to the index. This prevents an infinite loop if the last page is reached.
		for j := 0; j < save; j++ {
			itm := items[j].Itemize()
			fmt.Println(itm.ID)
			(*index)[itm.ID] = itm
		}

		// Prevent getting hit by the rate limiter.
		time.Sleep(2 * time.Second)
	}

	// Add names to our exported pages. Page names are just an incrementing sequence of integers, with the most recent page having the highest value.
	// Processing the slice backwards makes it so that the most recent page (the one that was first exported) is the last item in the pages slice.
	for i := len(exported) - 1; i >= 0; i-- {
		// Save the update files in their own directory.
		exported[i].Path = directory + "/exported/updated/" + strconv.Itoa(len(pages)+len(exported)-i) + ".json"

		// itemizers = append(itemizers, exported[i])
	}

	// Reorder all of our saved items.
	// Shuffle(itemizers, 30)

	// Write out all of our updated items.
	for i := range exported {
		err := exported[i].Write()
		if err != nil {
			fmt.Println("Error writing to file:", err)
		}
	}

	return nil
}
