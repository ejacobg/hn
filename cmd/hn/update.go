package main

import (
	"flag"
	"fmt"
	"github.com/ejacobg/hn/internal/args"
	"github.com/ejacobg/hn/internal/auth"
	"github.com/ejacobg/hn/internal/item"
	"github.com/ejacobg/hn/internal/update"
	"os"
)

var (
	updateFlags = flag.NewFlagSet("update", flag.ExitOnError)
	updateDir   = updateFlags.String("directory", "./<favorite|upvoted>/<submissions|comments>", "(Optional) Directory to be updated.")
	shuffle     = updateFlags.Bool("shuffle", false, "Shuffle items in the given directory.")
)

func init() {
	updateFlags.Usage = func() {
		w := updateFlags.Output()
		fmt.Fprintln(w, "Usage of update:")
		fmt.Fprintln(w, "update <favorite|upvoted> <submissions|comments> <username> [flags]")
		fmt.Fprintln(w, "update <favorite|upvoted> <submissions|comments> -shuffle")
		fmt.Fprintln(w, "update <-h|-help>")
		updateFlags.PrintDefaults()
		fmt.Fprintln(w, "To view upvoted posts, a password or token is required.")
	}
}

func updateCmd() {
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

	user.Username = os.Args[4]
	if user.Username == "-shuffle" {
		*shuffle = true
	}

	if *updateDir == "./<favorite|upvoted>/<submissions|comments>" {
		*updateDir = "./" + saveType + "/" + itemType
	}

	if *shuffle {
		switch itemType {
		case "submissions":
			update.ShuffleDir[item.Story](*updateDir, 30)
		case "comments":
			update.ShuffleDir[item.Comment](*updateDir, 30)
		}
		os.Exit(0)
	}

	var err error
	switch itemType {
	case "submissions":
		err = update.Items[item.Story](*updateDir, saveType, itemType, user)
	case "comments":
		err = update.Items[item.Comment](*updateDir, saveType, itemType, user)
	}

	if err != nil {
		fmt.Println("Error updating items:", err)
	}
}
