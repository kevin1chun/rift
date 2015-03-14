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
	logLevel := flags.String("log", "fatal", "")
	
	flags.Parse(os.Args[1:])

	logging.CurrentLevel = logging.ToLevel(*logLevel)

	args := flags.Args()

	switch {
	default:
		flags.Usage()
	case len(args) == 1 && args[0] == "version":
		printVersion()
	// case len(args) > 1 && args[0] == "build":
	// 	build(args[1:])
	case len(args) > 1 && args[0] == "run":
		run(args[1:])
	}
}

func printUsage() {
	fmt.Printf("Usage: rift [OPTIONS] COMMAND [ARGS]\n\n" +
		"OPTIONS\n" +
		"\t--log LEVEL\tSets the log level (default is \"FATAL\")\n" +
		"\n" +
		"COMMANDS\n" +
		"\tversion\t\tPrints the Rift version\n" +
		// "\tbuild\tBuilds the provided source files\n" +
		"\trun\t\tBuilds and runs the provided source files\n" +
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
