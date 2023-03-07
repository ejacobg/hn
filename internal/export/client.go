package export

import (
	"errors"
	"fmt"
	"github.com/ejacobg/hn/internal/auth"
	"github.com/ejacobg/hn/internal/item"
	"golang.org/x/net/html"
	"net/http"
)

type Client struct {
	user   auth.User
	client *http.Client
}

// NewClient returns a Client authenticated with the given user.
func NewClient(user auth.User) (*Client, error) {
	client := Client{user: user}

	var err error
	// If the token is missing, log in with the password to obtain it.
	if user.Token == "" {
		// If the password is missing, fail.
		if user.Password == "" {
			return nil, errors.New("NewClient: either token or password are required")
		}

		// Log in and obtain the token.
		client.client, err = auth.Login(user.Username, user.Password)
		if err != nil {
			return nil, err
		}
	} else {
		// Otherwise, authenticate with the token.
		client.client, err = auth.Token(user.Token)
		if err != nil {
			return nil, err
		}
	}

	return &client, nil
}

// DefaultClient returns a Client associated with (but not authenticated with) the given user.
func DefaultClient(user auth.User) *Client {
	return &Client{user, http.DefaultClient}
}

// GetPage will return a slice of items from the user's saved items.
func (c *Client) GetPage(saveType, itemType string, page int) ([]item.Itemizer, error) {
	doc, err := c.getDocument(saveType, itemType, page)
	if err != nil {
		return nil, err
	}

	switch itemType {
	case "submissions":
		items, err := Submissions(doc)
		if err != nil {
			return nil, err
		}
		return item.Itemizers(items), nil
	case "comments":
		items, err := Comments(doc)
		if err != nil {
			return nil, err
		}
		return item.Itemizers(items), nil
	}

	return nil, fmt.Errorf("GetPage: invalid item type %q", itemType)
}

// getDocument returns a parsed HTML document for the given page of items.
func (c *Client) getDocument(saveType, itemType string, page int) (*html.Node, error) {
	// Generate request.
	req, err := c.user.NewRequest(saveType, itemType, page)
	if err != nil {
		return nil, err
	}

	// 	Send request.
	resp, err := c.client.Get(req.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)

	// 	Parse request.
	return html.Parse(resp.Body)
}
