package sanity

import (
	"fmt"
)

func Fail(failure string, vars...interface{}) {
	panic(fmt.Errorf("Assertion failed: %s", fmt.Sprintf(failure, vars...)))
}

func Ensure(expr bool, otherwise string, vars...interface{}) {
	if !expr {
		Fail(otherwise, vars...)
	}
}