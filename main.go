package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	hn       = flag.NewFlagSet("hn", flag.ExitOnError)
	password = hn.String("password", "", "Password for the given user.")
	token    = hn.String("token", "", "Value of the 'user' cookie from a logged-in session. Takes priority over password.")
	page     = hn.Int("page", 1, "(Optional) Which page to read from.")
)

func main() {
	var saveType, postType, username string

	if len(os.Args) <= 4 {
		fmt.Println("Too few arguments.")
		os.Exit(1)
	}

	hn.Parse(os.Args[4:])
	args := os.Args[1:4]
	fmt.Println(args)

	switch saveType = args[0]; saveType {
	case "favorite":
		fmt.Println("select favorites")
	case "upvoted":
		fmt.Println("select upvoted")
	default:
		fmt.Println("unrecognized save type", saveType)
		os.Exit(1)
	}

	switch postType = args[1]; postType {
	case "submissions":
		fmt.Println("from submissions")
	case "comments":
		fmt.Println("from comments")
	default:
		fmt.Println("unrecognized post type", postType)
		os.Exit(1)
	}

	username = args[2]
	fmt.Println("where user =", username)

	// If the token is missing, log in with the password to obtain it.
	if *token == "" {
		// If the password is not given, fail.
		if *password == "" {
			fmt.Println("Either token or password are required.")
			os.Exit(1)
		}

		fmt.Println("with password", *password)

		// Log in and obtain the token.
		// ...
		// hn.Set("token", value)
	} else {
		fmt.Println("with token", *token)
	}
}
