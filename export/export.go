package export

import (
	"encoding/json"
	"github.com/ejacobg/hn/internal/scrape"
	"golang.org/x/net/html"
	"io"
)

// Submissions will parse the given document for stories, and write them as a JSON array to the destination.
func Submissions(doc *html.Node, dst io.Writer) error {
	posts := scrape.Submissions(doc)
	submissions, err := fromSubmissions(posts)
	if err != nil {
		return err
	}

	output, err := json.MarshalIndent(submissions, "", "\t")
	if err != nil {
		return err
	}

	_, err = dst.Write(output)
	return err
}

// Comments functions similarly to Submissions, excepts parses the document for comments.
func Comments(doc *html.Node, dst io.Writer) error {
	posts := scrape.Comments(doc)
	comments, err := fromComments(posts)
	if err != nil {
		return err
	}

	output, err := json.MarshalIndent(comments, "", "\t")
	if err != nil {
		return err
	}

	_, err = dst.Write(output)
	return err
}
