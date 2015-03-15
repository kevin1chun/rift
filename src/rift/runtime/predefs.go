package runtime

import (
	"fmt"
	"os"
	"strings"
	"rift/support/collections"
	"rift/support/logging"
)

// TODO: Type checks on arguments

var Predefs collections.PersistentMap

func InitPredefs() {
	Predefs = collections.NewPersistentMap()
	Predefs.Set("std:printf", printf)
	Predefs.Set("std:println", println)
	Predefs.Set("std:exit", exit)
	// Predefs.Set("std:file", file)

	logging.Debug("Built-in environment:")
	for k, _ := range Predefs.Freeze() {
		logging.Debug(" |- %s", k)
	}
}

func printf(args []interface{}) interface{} {
	fmt.Printf(args[0].(string), args[1:]...)
	return nil
}

func println(args []interface{}) interface{} {
	var stringedArgs []string
	for _, arg := range args {
		stringedArgs = append(stringedArgs, fmt.Sprintf("%v", arg))
	}
	fmt.Println(strings.Join(stringedArgs, ""))
	return nil
}

func exit(args []interface{}) interface{} {
	ensureArity("std:exit", 1, len(args))
	os.Exit(args[0].(int))
	return nil
}

// func file(args[]interface{}) interface{} {
// 	return nil
// }
