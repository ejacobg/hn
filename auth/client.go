package auth

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

type User struct {
	Username string
	Password string
	Token    string
}

// NewClient returns a *http.Client authenticated with the given user.
func (u User) NewClient() (*http.Client, error) {
	// If the token is missing, log in with the password to obtain it.
	if u.Token == "" {
		// If the password is not given, fail.
		if u.Password == "" {
			return nil, errors.New("NewClient: either token or password are required")
		}

		// Log in and obtain the token.
		return Login(u.Username, u.Password)
	} else {
		// Otherwise, authenticate with the token.
		return Token(u.Token)
	}
}

// NewRequest returns a request URL with query parameters set to the given values.
func (u User) NewRequest(saveType, itemType string, page int) (*url.URL, error) {
	// Create query parameters.
	query := url.Values{}
	query.Set("id", u.Username)
	query.Set("p", strconv.Itoa(page))
	if itemType == "comments" {
		query.Set("comments", "t")
	}

	// Create request URL.
	reqURL, err := url.Parse(BaseURL + "/" + saveType)
	if err != nil {
		return nil, err
	}

	// Attach query info to request URL.
	reqURL.RawQuery = query.Encode()

	return reqURL, nil
}
