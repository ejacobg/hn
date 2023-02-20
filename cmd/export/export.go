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
	export   = flag.NewFlagSet("export", flag.ExitOnError)
	password = export.String("password", "", "Password for the given user.")
	token    = export.String("token", "", "Value of the 'user' cookie from a logged-in session. Takes priority over password.")
	page     = export.Int("page", 1, "(Optional) Which page to read from.")
)

func init() {
	export.Usage = func() {
		w := export.Output()
		fmt.Fprintln(w, "Usage of export:")
		fmt.Fprintln(w, "export <favorite|upvoted> <submissions|comments> <username> [flags]")
		fmt.Fprintln(w, "export <-h|-help>")
		export.PrintDefaults()
		fmt.Fprintln(w, "To view upvoted posts, a password or token is required.")
	}
}

func main() {
	var saveType, itemType, username string
	var client *http.Client

	// export requires 3 arguments other than the command name.
	if len(os.Args) < 4 {
		// If there are less than 3 arguments, check if one of them is for the help flag.
		for i := 1; i < len(os.Args); i++ {
			if os.Args[i] == "-h" || os.Args[i] == "-help" {
				export.Usage()
				os.Exit(0)
			}
		}
		fmt.Println("Too few arguments.")
		export.Usage()
		os.Exit(1)
	}

	export.Parse(os.Args[4:])
	args := os.Args[1:4]
	// fmt.Println(args)

	switch saveType = args[0]; saveType {
	case "favorite":
		// fmt.Println("select favorites")
		saveType = "favorites"
	case "upvoted":
		// fmt.Println("select upvoted")
	default:
		fmt.Println("Unrecognized save type:", saveType)
		os.Exit(1)
	}

	switch itemType = args[1]; itemType {
	case "submissions":
		// fmt.Println("from submissions")
	case "comments":
		// fmt.Println("from comments")
	default:
		fmt.Println("Unrecognized item type:", itemType)
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
	if itemType == "comments" {
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