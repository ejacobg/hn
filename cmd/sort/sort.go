package main

import (
	"flag"
	"fmt"
	"github.com/ejacobg/hn/args"
	"github.com/ejacobg/hn/item"
	"github.com/ejacobg/hn/sort"
	"os"
)

var sortFlags = flag.NewFlagSet("sort", flag.ExitOnError)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Too few arguments.")
		sortFlags.Usage()
		os.Exit(2)
	}

	saveType, itemType, code := args.Parse(os.Args[2:], sortFlags.Usage)
	if code >= 0 {
		os.Exit(code)
	}

	var err error
	switch itemType {
	case "submissions":
		err = sort.Items[item.Story](saveType, itemType)
	case "comments":
		err = sort.Items[item.Comment](saveType, itemType)
	}

	if err != nil {
		fmt.Println("Failed to sort items:", err)
		os.Exit(1)
	}
}
