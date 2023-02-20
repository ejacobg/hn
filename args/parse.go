package args

import "fmt"

// Parse will process an argument slice of the form: <favorite|upvoted> <submissions|comments>
// If the arguments do not conform to the above format, then the usage function will be used to print instructions to the screen, and code 2 will be returned.
// If one of the arguments is -h or -help, then the usage instructions will be printed, and code 0 will be returned.
// Otherwise, the correct save type and item type will be returned, along with code -1.
func Parse(args []string, usage func()) (saveType, itemType string, code int) {
	if len(args) < 2 {
		for i := 0; i < len(args); i++ {
			if args[i] == "-h" || args[i] == "-help" {
				usage()
				return "", "", 0
			}
		}
		fmt.Println("Too few arguments.")
		usage()
		return "", "", 2
	}

	saveType, itemType = args[1], args[2]

	if !(saveType == "favorite" || saveType == "upvoted") {
		fmt.Println("Unrecognized save type:", saveType)
		code = 2
		return
	}

	if !(itemType == "submissions" || itemType == "comments") {
		fmt.Println("Unrecognized item type:", itemType)
		code = 2
		return
	}

	code = -1
	return
}
