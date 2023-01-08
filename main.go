package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ejacobg/hn/auth"
	"github.com/ejacobg/hn/item"
	"github.com/ejacobg/hn/scrape"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

var (
	hn       = flag.NewFlagSet("hn", flag.ExitOnError)
	password = hn.String("password", "", "Password for the given user.")
	token    = hn.String("token", "", "Value of the 'user' cookie from a logged-in session. Takes priority over password.")
	page     = hn.Int("page", 1, "(Optional) Which page to read from.")
)

func init() {
	hn.Usage = func() {
		w := hn.Output()
		fmt.Fprintln(w, "Usage of hn:")
		fmt.Fprintln(w, "hn <favorite|upvoted> <submissions|comments> <username> [flags]")
		fmt.Fprintln(w, "hn <-h|-help>")
		hn.PrintDefaults()
		fmt.Fprintln(w, "To view upvoted posts, a password or token is required.")
	}
}

func main() {
	var saveType, postType, username string
	var client *http.Client

	// hn requires 3 arguments other than the command name.
	if len(os.Args) < 4 {
		// If there are less than 3 arguments, check if one of them is for the help flag.
		for i := 1; i < len(os.Args); i++ {
			if os.Args[i] == "-h" || os.Args[i] == "-help" {
				hn.Usage()
				os.Exit(0)
			}
		}
		fmt.Println("Too few arguments.")
		hn.Usage()
		os.Exit(1)
	}

	hn.Parse(os.Args[4:])
	args := os.Args[1:4]
	fmt.Println(args)

	switch saveType = args[0]; saveType {
	case "favorite":
		// fmt.Println("select favorites")
		saveType = "favorites"
	case "upvoted":
		// fmt.Println("select upvoted")
	default:
		fmt.Println("unrecognized save type", saveType)
		os.Exit(1)
	}

	switch postType = args[1]; postType {
	case "submissions":
		// fmt.Println("from submissions")
	case "comments":
		// fmt.Println("from comments")
	default:
		fmt.Println("unrecognized post type", postType)
		os.Exit(1)
	}

	username = args[2]
	// fmt.Println("where user =", username)

	// Upvoted posts require an authenticated client.
	if saveType == "upvoted" {
		// If the token is missing, log in with the password to obtain it.
		if *token == "" {
			// If the password is not given, fail.
			if *password == "" {
				fmt.Println("Either token or password are required.")
				os.Exit(1)
			}

			// fmt.Println("with password", *password)

			// Log in and obtain the token.
			var err error
			client, err = auth.Login(username, *password)
			if err != nil {
				fmt.Println("Unable to log in:", err)
				os.Exit(1)
			}
		} else {
			// fmt.Println("with token", *token)

			var err error
			client, err = auth.Token(*token)
			if err != nil {
				fmt.Println("Unable to set token:", err)
				os.Exit(1)
			}
		}
	} else {
		client = http.DefaultClient
	}

	query := url.Values{}
	query.Set("id", username)
	query.Set("p", strconv.Itoa(*page))
	if postType == "comments" {
		query.Set("comments", "t")
	}

	reqURL, err := url.Parse(auth.BaseURL + "/" + saveType)
	if err != nil {
		fmt.Println("Failed to parse URL:", err)
		os.Exit(1)
	}
	reqURL.RawQuery = query.Encode()

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
	switch postType {
	case "submissions":
		posts = scrape.Submissions(doc)
	case "comments":
		posts = scrape.Comments(doc)
	}

	for _, post := range posts {
		switch postType {
		case "submissions":
			submission, err := item.FromSubmission(post)
			if err != nil {
				fmt.Println("Error parsing HTML:", err)
				os.Exit(1)
			}
			output, err := json.MarshalIndent(submission, "", "\t")
			if err != nil {
				fmt.Println("Error marshaling data:", err)
				os.Exit(1)
			}
			fmt.Println(string(output) + ",")
		case "comments":
			comment, err := item.FromComment(post)
			if err != nil {
				fmt.Println("Error parsing HTML:", err)
				os.Exit(1)
			}
			output, err := json.MarshalIndent(comment, "", "\t")
			if err != nil {
				fmt.Println("Error marshaling data:", err)
				os.Exit(1)
			}
			fmt.Println(string(output) + ",")
		}
	}
}
