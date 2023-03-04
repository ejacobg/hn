package auth

import (
	"net/url"
	"testing"
)

func TestLogin(t *testing.T) {
	client, err := Login("<username>", "<password>")
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}

	URL, err := url.Parse(BaseURL)
	if err != nil {
		t.Fatalf("could not parse URL: %v", err)
	}

	cookies := client.Jar.Cookies(URL)
	var found bool
	for _, c := range cookies {
		t.Log(c.String())
		if c.Name == "user" {
			found = true
		}
	}
	if !found {
		t.Log("no 'user' cookie found")
		t.Fail()
	}
}
