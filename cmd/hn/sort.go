package main

import (
	"flag"
	"fmt"
	"github.com/ejacobg/hn/internal/args"
	"github.com/ejacobg/hn/internal/item"
	"github.com/ejacobg/hn/internal/sort"
	"os"
)

var (
	sortFlags = flag.NewFlagSet("sort", flag.ExitOnError)
	sortDir   = sortFlags.String("directory", "", "(Optional) Custom directory to read from. Must be structured correctly.")
)

func init() {
	sortFlags.Usage = func() {
		w := sortFlags.Output()
		fmt.Fprintln(w, "Usage of sort:")
		fmt.Fprintln(w, "sort <favorite|upvoted> <submissions|comments> [flags]")
		fmt.Fprintln(w, "sort <-h|-help>")
		sortFlags.PrintDefaults()
	}
}

func sortCmd() {
	if len(os.Args) < 4 {
		fmt.Println("Too few arguments.")
		sortFlags.Usage()
		os.Exit(2)
	}

	saveType, itemType, code := args.Parse(os.Args[2:], sortFlags.Usage)
	if code >= 0 {
		os.Exit(code)
	}
	sortFlags.Parse(os.Args[4:])

	if *sortDir == "" {
		*sortDir = "./" + saveType + "/" + itemType
	}

	var err error
	switch itemType {
	case "submissions":
		err = sort.Items[item.Story](*sortDir)
	case "comments":
		err = sort.Items[item.Comment](*sortDir)
	}

	if err != nil {
		fmt.Println("Failed to sort items:", err)
		os.Exit(1)
	}
}
