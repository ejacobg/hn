package main

import (
	"flag"
	"fmt"
	"github.com/ejacobg/hn/args"
	"github.com/ejacobg/hn/auth"
	"os"
)

var (
	updateFlags = flag.NewFlagSet("update", flag.ExitOnError)
	directory   = updateFlags.String("directory", "./<favorite|upvoted>/<submissions|comments>", "(Optional) Directory to be updated.")
)

func init() {
	updateFlags.Usage = func() {
		w := updateFlags.Output()
		fmt.Fprintln(w, "Usage of update:")
		fmt.Fprintln(w, "update <favorite|upvoted> <submissions|comments> <username> [flags]")
		fmt.Fprintln(w, "update <-h|-help>")
		updateFlags.PrintDefaults()
		fmt.Fprintln(w, "To view upvoted posts, a password or token is required.")
	}
}

func main() {
	var user auth.User
	updateFlags.StringVar(&user.Password, "password", "", "Password for the given user.")
	updateFlags.StringVar(&user.Token, "token", "", "Value of the 'user' cookie from a logged-in session. Takes priority over password.")

	if len(os.Args) < 5 {
		fmt.Println("Too few arguments.")
		updateFlags.Usage()
		os.Exit(2)
	}

	saveType, itemType, code := args.Parse(os.Args[2:], updateFlags.Usage)
	if code >= 0 {
		os.Exit(code)
	}
	updateFlags.Parse(os.Args[5:])

	if saveType == "favorite" {
		// URL route uses "favorites"
		saveType = "favorites"
	}

	user.Username = os.Args[4]

	// var client *http.Client
	//
	// // Upvoted posts require an authenticated client.
	// if saveType == "upvoted" {
	// 	var err error
	// 	client, err = user.NewClient()
	// 	if err != nil {
	// 		fmt.Println("Error creating client:", err)
	// 	}
	// } else {
	// 	client = http.DefaultClient
	// }

	fmt.Println("Updating", saveType, itemType)
}
