package main

import (
	"flag"
	"fmt"
	"os"
	"rift/lang"
	"rift/vm"
)

func main() {
	flag.Parse()
	args := flag.Args()

	switch {
	default:
		printUsage()
	case len(args) == 1 && args[0] == "version":
		printVersion()
	case len(args) > 1 && args[0] == "build":
		build(args[1:])
	case len(args) > 1 && args[0] == "run":
		run(args[1:])
	}
}

func printUsage() {
	fmt.Printf("Usage: rift version|(run <filename>)\n")
	os.Exit(1)
}

func printVersion() {
	fmt.Println("rift-v0.1")
}

func build(filenames []string) []*lang.Node {
	var rifts []*lang.Node
	for _, filename := range filenames {
		source, readErr := os.Open(filename)
		if readErr != nil {
			fmt.Printf("Couldn't open file [%s]: %+v", filename, readErr)
		}

		parsed, parseErr := lang.Parse(source)
		if parseErr != nil {
			fmt.Printf("Parse error: %+v", parseErr)
		}

		fmt.Printf("Compiled file [%s]\n", filename)

		rifts = append(rifts, parsed.Rifts()...)
	}

	return rifts
}

func run(filenames []string) {
	vm := vm.New(build(filenames))
	vm.Run()
}