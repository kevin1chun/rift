package logging

import (
	"fmt"
)

type Level int

const (
	DEBUG = iota
	INFO
	WARN
	FATAL
)

// TODO: Thread-safety
var CurrentLevel Level = DEBUG

func log(level Level, levelName string, msg string, args...interface{}) {
	if level >= CurrentLevel {
		fmt.Printf("[%s] %s\n", levelName, fmt.Sprintf(msg, args))
	}
}

func Debug(msg string, args...interface{}) {
	log(DEBUG, "debug", msg, args...)
}

func Info(msg string, args...interface{}) {
	log(INFO, "info", msg, args...)
}

func Warn(msg string, args...interface{}) {
	log(WARN, "warn", msg, args...)
}

func Fatal(msg string, args...interface{}) {
	log(FATAL, "fatal", msg, args...)
}