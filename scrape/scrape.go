package scrape

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Note: the structure of HN's favorites/upvoted pages may change in the future, affecting the function of this code.

// Submissions and Comments happen to use the same algorithm.

func Submissions(doc *html.Node) []*html.Node {
	matcher := func(node *html.Node) (bool, bool) {
		if node.Type == html.ElementNode && node.DataAtom == atom.Tr {
			for _, a := range node.Attr {
				if a.Key == "class" && a.Val == "athing" {
					return true, true
				}
			}
		}
		return false, false
	}
	return FindNodes(doc, matcher)
}

func Comments(doc *html.Node) []*html.Node {
	matcher := func(node *html.Node) (bool, bool) {
		if node.Type == html.ElementNode && node.DataAtom == atom.Tr {
			for _, a := range node.Attr {
				if a.Key == "class" && a.Val == "athing" {
					return true, true
				}
			}
		}
		return false, false
	}
	return FindNodes(doc, matcher)
}

// GetElementWithClass returns the first element underneath and including `node` that has the
// given class value (as given in the HTML). The classes must be in the same order as those given.
func GetElementWithClass(node *html.Node, tagname atom.Atom, class string) *html.Node {
	matcher := func(node *html.Node) bool {
		if node.Type == html.ElementNode && node.DataAtom == tagname {
			for _, a := range node.Attr {
				if a.Key == "class" && a.Val == class {
					return true
				}
			}
		}
		return false
	}
	return FindNode(node, matcher)
}

// FindNode will return the first node that the matcher function accepts.
func FindNode(node *html.Node, matcher func(node *html.Node) bool) *html.Node {
	if matcher(node) {
		return node
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		n := FindNode(c, matcher)
		if n != nil {
			return n
		}
	}
	return nil
}

// FindNodes will gather all nodes that the matcher function accepts.
// The matcher function's second return value indicates whether to search the children of the current node.
func FindNodes(node *html.Node, matcher func(node *html.Node) (bool, bool)) (nodes []*html.Node) {
	var keep, exit bool
	var f func(*html.Node)
	f = func(n *html.Node) {
		keep, exit = matcher(n)
		if keep {
			nodes = append(nodes, n)
		}
		if exit {
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(node)
	return nodes
}
