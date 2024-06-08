package export

import (
	"errors"

	"github.com/ejacobg/hn/internal/auth"
	"github.com/ejacobg/hn/internal/item"
	"github.com/ejacobg/hn/internal/scrape"
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

	athingTr := scrape.GetElementWithClass(node, atom.Tr, "athing")
	if athingTr == nil {
		return comment, errors.New("fromComment: could not find tr.athing")
	}

	onstorySpan := scrape.GetElementWithClass(athingTr, atom.Span, "onstory")
	if onstorySpan == nil {
		return comment, errors.New("fromComment: could not find tr.athing span.onstory")
	}

	var a *html.Node
	for a = onstorySpan.FirstChild; a != nil; a = a.NextSibling {
		if a.Type == html.ElementNode && a.DataAtom == atom.A {
			break
		}
	}
	if a == nil {
		return comment, errors.New("fromComment: could not find tr.athing span.onstory > a")
	}

	commentDiv := scrape.GetElementWithClass(athingTr, atom.Div, "comment")
	if commentDiv == nil {
		return comment, errors.New("fromComment: could not find tr.athing div.comment")
	}

	var div *html.Node
	for div = commentDiv.FirstChild; div != nil; div = div.NextSibling {
		if div.Type == html.ElementNode && div.DataAtom == atom.Div {
			break
		}
	}
	if div == nil {
		return comment, errors.New("fromComment: could not find tr.athing div.comment > div")
	}

	// Grab all text nodes.
	// The HTML returned from the server DOES NOT close its <p> tags, which results in an extra text node and whitespace being created after parsing.
	// I will leave these artifacts alone since they are inconsequential to the results of the program. The tests will account for these errors though.
	text := scrape.FindNodes(div, func(node *html.Node) (bool, bool) {
		return node.Type == html.TextNode, false
	})

	for _, attr := range athingTr.Attr {
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
