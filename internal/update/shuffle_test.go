package update

import (
	"github.com/ejacobg/hn/internal/item"
	"golang.org/x/exp/slices"
	"testing"
)

// testItem implements the item.Itemizer interface.
type testItem struct {
	*item.Item
	value int
}

func (ti testItem) Itemize() *item.Item {
	return ti.Item
}

func newPage(values ...int) (page item.Page[testItem]) {
	page.Items = []testItem{}
	for _, value := range values {
		page.Items = append(page.Items, testItem{nil, value})
	}
	return
}

func lengths(pages []item.Page[testItem]) (lens []int) {
	for _, page := range pages {
		lens = append(lens, len(page.Items))
	}
	return
}

func log(t *testing.T, pages []item.Page[testItem]) {
	for _, page := range pages {
		t.Log(page.Items)
	}
}

func TestShuffle(t *testing.T) {
	tests := []struct {
		name    string
		pages   []item.Page[testItem]
		lengths []int // Represents the lengths of each page after the shuffle.
	}{
		{"Unchanged", []item.Page[testItem]{newPage(1, 2, 3, 4, 5), newPage(6, 7, 8, 9, 10), newPage()}, []int{5, 5, 0}},
		{"First page filled", []item.Page[testItem]{newPage(), newPage(6, 7, 8, 9, 10), newPage()}, []int{5, 0, 0}},
		{"Last page filled", []item.Page[testItem]{newPage(1, 2, 3, 4, 5, 6, 7, 8, 9, 10), newPage(11, 12, 13, 14, 15), newPage()}, []int{5, 5, 5}},
		{"Last page overfilled", []item.Page[testItem]{newPage(1, 2, 3, 4, 5, 6, 7, 8, 9, 10), newPage(11, 12, 13, 14, 15), newPage(16, 17, 18, 19, 20)}, []int{5, 5, 10}},
	}

	limit := 5
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			log(t, test.pages)
			Shuffle(test.pages, limit)
			log(t, test.pages)
			lens := lengths(test.pages)
			if !slices.Equal(lens, test.lengths) {
				t.Errorf("got %v, want %v", lens, test.lengths)
			}
		})
	}
}
