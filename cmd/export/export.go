package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ejacobg/hn/args"
	"github.com/ejacobg/hn/auth"
	"github.com/ejacobg/hn/export"
	"golang.org/x/net/html"
	"io"
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

	// Upvoted posts require an authenticated client.
	var client *http.Client
	if saveType == "upvoted" {
		var err error
		client, err = user.NewClient()
		if err != nil {
			fmt.Println("Error creating client:", err)
		}
	} else {
		client = http.DefaultClient
	}

	// Generate request.
	reqURL, err := user.NewRequest(saveType, itemType, *page)
	if err != nil {
		fmt.Println("Failed to create URL:", err)
		os.Exit(1)
	}

	// Send request.
	resp, err := client.Get(reqURL.String())
	if err != nil {
		fmt.Println("Error retrieving content:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Parse request.
	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println("Error parsing document:", err)
		os.Exit(1)
	}

	var items any
	switch itemType {
	case "submissions":
		items, err = export.Submissions(doc)
	case "comments":
		items, err = export.Comments(doc)
	}

	if err != nil {
		fmt.Println("Error exporting items:", err)
		os.Exit(1)
	}

	write(items, os.Stdout)
}

func write(data any, dst io.Writer) error {
	output, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	_, err = dst.Write(output)
	return err
}
