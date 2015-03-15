package main

import (
	"flag"
	"fmt"
	"os"
	"rift/lang"
	"rift/support/logging"
	"rift/runtime"
)

const (
	RIFT_VERSION = "v0.1-alpha"
	INVALID_ARGS = 0
	INVALID_FILE = 1
	SYNTAX_ERROR = 2
)

func main() {
	flags := flag.NewFlagSet("rift", flag.ExitOnError)
	flags.Usage = printUsage

	showVersion  := flags.Bool("version", false, "Prints this version of Rift")
	debug := flags.Bool("verbose", false, "")
	
	flags.Parse(os.Args[1:])

	if *debug {
		logging.CurrentLevel = logging.DEBUG
	}

	args := flags.Args()

	switch {
	default:
		flags.Usage()
	case *showVersion:
		printVersion()
	case len(args) >= 1:
		run(args)
	}
}

func printUsage() {
	fmt.Printf("Usage: rift [OPTIONS] [FILES]\n\n" +
		"OPTIONS\n" +
		"  --verbose Prints verbose Rift interpreter logs\n" +
		"  --version Prints the Rift version\n" +
		"\n")
	os.Exit(INVALID_ARGS)
}

func printVersion() {
	fmt.Printf("rift-%s\n", RIFT_VERSION)
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
	runtime.Run(rifts)
	// initialCtx := runtime.BuildContext(rifts)
	// dispatcher := runtime.LocalDispatcher{}
	// vm := runtime.NewVM(&dispatcher)
	// vm.SwitchTo(initialCtx)
	// vm.Run()
}
