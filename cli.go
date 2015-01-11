package main

import (
	"flag"
	"fmt"
	"os"
	"rift/lang"
)

func main() {
	flag.Parse()
	filename := flag.Arg(0)
	source, readErr := os.Open(filename)
	if readErr != nil {
		fmt.Printf("Couldn't open file [%s]: %+v", filename, readErr)
	}

	parsed, parseErr := lang.Parse(source)
	if parseErr != nil {
		fmt.Printf("Parse error: %+v", parseErr)
	}

	fmt.Println(parsed.Lisp())
}