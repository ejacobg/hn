package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ejacobg/hn/args"
	"github.com/ejacobg/hn/auth"
	"github.com/ejacobg/hn/item"
	"github.com/ejacobg/hn/scrape"
	"golang.org/x/net/html"
	"net/http"
	"os"
)

var (
	exportFlags = flag.NewFlagSet("export", flag.ExitOnError)
	page        = exportFlags.Int("page", 1, "(Optional) Which page to read from.")
)

func init() {
	exportFlags.Usage = func() {
		w := exportFlags.Output()
		fmt.Fprintln(w, "Usage of export:")
		fmt.Fprintln(w, "export <favorite|upvoted> <submissions|comments> <username> [flags]")
		fmt.Fprintln(w, "export <-h|-help>")
		exportFlags.PrintDefaults()
		fmt.Fprintln(w, "To view upvoted posts, a password or token is required.")
	}
}

func main() {
	var user auth.User
	exportFlags.StringVar(&user.Password, "password", "", "Password for the given user.")
	exportFlags.StringVar(&user.Token, "token", "", "Value of the 'user' cookie from a logged-in session. Takes priority over password.")

	// fmt.Println(os.Args)

	if len(os.Args) < 5 {
		fmt.Println("Too few arguments.")
		exportFlags.Usage()
		os.Exit(1)
	}

	saveType, itemType, code := args.Parse(os.Args[2:], exportFlags.Usage)
	if code >= 0 {
		os.Exit(code)
	}
	exportFlags.Parse(os.Args[5:])

	if saveType == "favorite" {
		// URL route uses "favorites"
		saveType = "favorites"
	}

	user.Username = os.Args[4]

	var client *http.Client

	// Upvoted posts require an authenticated client.
	if saveType == "upvoted" {
		var err error
		client, err = user.NewClient()
		if err != nil {
			fmt.Println("Error creating client:", err)
		}
	} else {
		client = http.DefaultClient
	}

	reqURL, err := user.NewRequest(saveType, itemType, *page)
	if err != nil {
		fmt.Println("Failed to create URL:", err)
		os.Exit(1)
	}

	// fmt.Println(reqURL.String())
	resp, err := client.Get(reqURL.String())
	if err != nil {
		fmt.Println("Error retrieving content:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// fmt.Println(resp.Header.Get("Content-Type"))
	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println("Error parsing document:", err)
		os.Exit(1)
	}

	var posts []*html.Node
	switch itemType {
	case "submissions":
		posts = scrape.Submissions(doc)
		submissions, err := item.FromSubmissions(posts)
		if err != nil {
			fmt.Println("Error parsing HTML:", err)
		}
		output, err := json.MarshalIndent(submissions, "", "\t")
		if err != nil {
			fmt.Println("Error marshaling data:", err)
			os.Exit(1)
		}
		fmt.Println(string(output))
	case "comments":
		posts = scrape.Comments(doc)
		comments, err := item.FromComments(posts)
		if err != nil {
			fmt.Println("Error parsing HTML:", err)
		}
		output, err := json.MarshalIndent(comments, "", "\t")
		if err != nil {
			fmt.Println("Error marshaling data:", err)
			os.Exit(1)
		}
		fmt.Println(string(output))
	}
}
