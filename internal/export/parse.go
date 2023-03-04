package export

import (
	"errors"
	"github.com/ejacobg/hn/internal/auth"
	"github.com/ejacobg/hn/internal/scrape"
	"github.com/ejacobg/hn/item"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Note: the structure of HN's favorites/upvoted pages may change in the future, affecting the function of this code.

func fromSubmission(node *html.Node) (item.Story, error) {
	story := item.Story{Item: &item.Item{}}

	tr := scrape.GetElementWithClass(node, atom.Tr, "athing")
	if tr == nil {
		return story, errors.New("fromSubmission: could not find tr.athing")
	}

	span := scrape.GetElementWithClass(tr, atom.Span, "titleline")
	if span == nil {
		return story, errors.New("fromSubmission: could not find tr.athing span.titleline")
	}

	var a *html.Node
	for a = span.FirstChild; a != nil; a = a.NextSibling {
		if a.Type == html.ElementNode && a.DataAtom == atom.A {
			break
		}
	}
	if a == nil {
		return story, errors.New("fromSubmission: could not find tr.athing span.titleline > a")
	}

	for _, attr := range tr.Attr {
		if attr.Key == "id" {
			story.ID = attr.Val
		}
	}
	story.Title = a.FirstChild.Data
	for _, attr := range a.Attr {
		if attr.Key == "href" {
			story.URL = attr.Val
		}
	}
	story.Discussion = auth.BaseURL + "/item?id=" + story.ID

	return story, nil
}

// fromSubmissions returns all parsed nodes, successful or not, and returns the last error encountered.
func fromSubmissions(nodes []*html.Node) (submissions []item.Story, err error) {
	var submission item.Story
	for _, node := range nodes {
		submission, err = fromSubmission(node)
		submissions = append(submissions, submission)
	}
	return
}

func fromComment(node *html.Node) (item.Comment, error) {
	comment := item.Comment{Item: &item.Item{}}

	tr := scrape.GetElementWithClass(node, atom.Tr, "athing")
	if tr == nil {
		return comment, errors.New("fromComment: could not find tr.athing")
	}

	story := scrape.GetElementWithClass(tr, atom.Span, "onstory")
	if story == nil {
		return comment, errors.New("fromComment: could not find tr.athing span.onstory")
	}

	var a *html.Node
	for a = story.FirstChild; a != nil; a = a.NextSibling {
		if a.Type == html.ElementNode && a.DataAtom == atom.A {
			break
		}
	}
	if a == nil {
		return comment, errors.New("fromComment: could not find tr.athing span.onstory > a")
	}

	div := scrape.GetElementWithClass(tr, atom.Div, "comment")
	if div == nil {
		return comment, errors.New("fromComment: could not find tr.athing div.comment")
	}

	var span *html.Node
	for span = div.FirstChild; span != nil; span = span.NextSibling {
		if span.Type == html.ElementNode && span.DataAtom == atom.Span {
			break
		}
	}
	if span == nil {
		return comment, errors.New("fromComment: could not find tr.athing div.comment > span")
	}

	// Grab all text nodes.
	// The HTML returned from the server DOES NOT close its <p> tags, which results in an extra text node and whitespace being created after parsing.
	// I will leave these artifacts alone since they are inconsequential to the results of the program. The tests will account for these errors though.
	text := scrape.FindNodes(span, func(node *html.Node) (bool, bool) {
		return node.Type == html.TextNode, false
	})

	for _, attr := range tr.Attr {
		if attr.Key == "id" {
			comment.ID = attr.Val
		}
	}
	if a.FirstChild != nil {
		comment.Story = a.FirstChild.Data
	}
	for _, t := range text {
		comment.Text = append(comment.Text, t.Data)
	}
	comment.Context = auth.BaseURL + "/context?id=" + comment.ID
	for _, attr := range a.Attr {
		if attr.Key == "href" {
			comment.Discussion = auth.BaseURL + "/" + attr.Val
		}
	}

	return comment, nil
}

// fromComments returns all parsed nodes, successful or not, and returns the last error encountered.
func fromComments(nodes []*html.Node) (comments []item.Comment, err error) {
	var comment item.Comment
	for _, node := range nodes {
		comment, err = fromComment(node)
		comments = append(comments, comment)
	}
	return
}
