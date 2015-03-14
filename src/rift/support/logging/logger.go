package logging

import (
	"fmt"
	"strings"
)

type Level int

const (
	DEBUG = iota
	INFO
	WARN
	FATAL
)

// TODO: Thread-safety?
var CurrentLevel Level = FATAL

func ToLevel(levelString string) Level {
	switch strings.ToLower(levelString) {
	default:
		return CurrentLevel
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "warn":
		return WARN
	case "fatal":
		return FATAL
	}
}

func log(level Level, levelName string, msg string, args...interface{}) {
	if level >= CurrentLevel {
		fmt.Printf("[%s] %s\n", strings.ToUpper(levelName), fmt.Sprintf(msg, args...))
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