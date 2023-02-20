package main

import (
	"fmt"
	"github.com/ejacobg/hn/item"
)

func main() {
	fmt.Printf("%T %v\n", item.Advice, item.Advice)
	fmt.Printf("%T %v\n", item.Hobby, item.Hobby)
}
