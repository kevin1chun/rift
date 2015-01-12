package main

import (
	"flag"
	"fmt"
	"os"
	"rift/lang"
	"rift/runtime"
)

const (
	INVALID_ARGS = 0
	INVALID_FILE = 1
	SYNTAX_ERROR = 2
)

func main() {
	flag.Parse()
	args := flag.Args()

	switch {
	default:
		printUsage()
	case len(args) == 1 && args[0] == "version":
		printVersion()
	// case len(args) > 1 && args[0] == "build":
	// 	build(args[1:])
	case len(args) > 1 && args[0] == "run":
		run(args[1:])
	}
}

func printUsage() {
	fmt.Printf("Usage: rift COMMAND [ARGS]\n\n" +
		"COMMANDS\n" +
		"\tversion\tPrints the Rift version\n" +
		// "\tbuild\tBuilds the provided source files\n" +
		"\trun\tBuilds and runs the provided source files\n" +
		"\n")
	os.Exit(INVALID_ARGS)
}

func printVersion() {
	fmt.Println("rift-v0.1")
}

func build(filenames []string) []*lang.Node {
	var rifts []*lang.Node
	for _, filename := range filenames {
		source, readErr := os.Open(filename)
		if readErr != nil {
			fmt.Printf("Couldn't open file [%s]: %+v\n", filename, readErr)
			os.Exit(INVALID_FILE)
		}

		parsed, parseErr := lang.Parse(source)
		if parseErr != nil {
			fmt.Printf("Syntax error [%s]: %+v\n", filename, parseErr)
			os.Exit(SYNTAX_ERROR)
		}

		rifts = append(rifts, parsed.Rifts()...)
	}

	return rifts
}

func run(filenames []string) {
	rifts := build(filenames)
	dispatcher := runtime.NewLocalDispatcher(rifts)
	vm := runtime.NewVM(dispatcher)
	vm.Run()
}