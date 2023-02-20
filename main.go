package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "-h", "-help":
		usage()
	case "export":
		fmt.Println("Exporting items...")
	case "sort":
		fmt.Println("Sorting items...")
	case "update":
		fmt.Println("Updating items...")
	default:
		fmt.Println("Unrecognized subcommand:", os.Args[1])
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Println("Usage of hn:")
	fmt.Println("\thn <export|sort|update> <favorite|upvoted> <submissions|comments> [args...] [flags...]")
	fmt.Println("\thn <-h|-help>")
	fmt.Println("Subcommands:")
	fmt.Println("\thn export - returns a user's saved items as JSON")
	fmt.Println("\thn sort - applies categories to exported items")
	fmt.Println("\thn update - organizes export files, pulls unexported items")
}
