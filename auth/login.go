package auth

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

const BaseURL = "https://news.ycombinator.com"

// Login returns a *http.Client populated with the `user` cookie set upon successful login.
func Login(username, password string) (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("login: failed to create cookiejar: %w", err)
	}

	client := http.Client{
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			// A successful login returns a redirect. We don't need to follow it.
			return http.ErrUseLastResponse
		},
		Jar: jar,
	}

	creds := url.Values{}
	creds.Set("acct", username)
	creds.Set("pw", password)

	resp, err := client.PostForm(BaseURL+"/login", creds)
	if err != nil {
		return nil, fmt.Errorf("login: login attempt failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 302 {
		return nil, fmt.Errorf("login: expected status code 302 but got %d", resp.StatusCode)
	}

	return &client, nil
}

// Token returns a *http.Client with a manually-set `user` cookie.
func Token(token string) (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("token: failed to create cookiejar: %w", err)
	}

	URL, err := url.Parse(BaseURL)
	if err != nil {
		return nil, fmt.Errorf("token: could not parse URL: %w", err)
	}

	cookie := http.Cookie{
		Name:  "user",
		Value: token,
	}
	jar.SetCookies(URL, []*http.Cookie{&cookie})

	client := http.Client{
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			// A successful login returns a redirect. We don't need to follow it.
			return http.ErrUseLastResponse
		},
		Jar: jar,
	}

	return &client, nil
}
