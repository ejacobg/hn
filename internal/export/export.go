package export

import (
	"github.com/ejacobg/hn/internal/item"
	"github.com/ejacobg/hn/internal/scrape"
	"golang.org/x/net/html"
)

// Submissions will parse the given document for stories, returning the stories it has found.
func Submissions(doc *html.Node) ([]item.Story, error) {
	posts := scrape.Submissions(doc)
	return fromSubmissions(posts)
}

// Comments functions similarly to Submissions, excepts parses the document for comments.
func Comments(doc *html.Node) ([]item.Comment, error) {
	posts := scrape.Comments(doc)
	return fromComments(posts)
}
