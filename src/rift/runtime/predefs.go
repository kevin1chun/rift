package runtime

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"rift/support/collections"
	"rift/support/logging"
)

var Predefs collections.PersistentMap

func InitPredefs() {
	Predefs = collections.NewPersistentMap()
	Predefs.Set("std:len", length)
	Predefs.Set("std:sprintf", sprintf)
	Predefs.Set("std:printf", printf)
	Predefs.Set("std:println", println)
	Predefs.Set("std:exit", exit)
	Predefs.Set("std:open", fileOpen)
	Predefs.Set("std:write", fileWrite)
	Predefs.Set("std:close", fileClose)

	logging.Debug("Built-in environment:")
	for k, _ := range Predefs.Freeze() {
		logging.Debug(" |- %s", k)
	}
}

func length(args []interface{}) interface{} {
	ensureArity(1, len(args))
	return len(args[0].(string))
}

func sprintf(args []interface{}) interface{} {
	return fmt.Sprintf(args[0].(string), args[1:]...)
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
	ensureArity(1, len(args))
	os.Exit(args[0].(int))
	return nil
}

func fileOpen(args []interface{}) interface{} {
	ensureArity(1, len(args))
	filename := args[0].(string)
	// TODO: Pass in mode as string?
	fd, _ := syscall.Open(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	return fd
}

func fileWrite(args []interface{}) interface{} {
	ensureArity(2, len(args))
	fd := args[0].(int)
	data := args[1].(string)
	syscall.Write(fd, []byte(data))
	return nil
}

func fileClose(args []interface{}) interface{} {
	ensureArity(1, len(args))
	fd := args[0].(int)
	syscall.Close(fd)
	return nil
}
