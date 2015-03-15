package runtime

import (
	"rift/support/sanity"
)

func ensureArity(refStr string, expectedLength int, actualLength int) {
	sanity.Ensure(actualLength == expectedLength, "Function [%s] expects [%d] arguments, but got [%d]", refStr, expectedLength, actualLength)
}