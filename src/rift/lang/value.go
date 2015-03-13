package lang

import (
	"strconv"
	"rift/support/sanity"
)

// TODO: Actually, is there a difference between a string node
// vs. the string value representing a named ref?
func (stringNode *Node) Str() string {
	sanity.Ensure(stringNode.Type == STRING, "Invalid cast from type [%s] to [%s]", stringNode.Type, STRING)
	return stringNode.Values[0].(string)
}

func (numericNode *Node) Int() int {
	sanity.Ensure(numericNode.Type == NUM, "Invalid cast from type [%s] to [%s]", numericNode.Type, NUM)
	intAsString := numericNode.Values[0].(string)
	intValue, _ := strconv.Atoi(intAsString)
	// TODO: WTF?
	// intValue, parseErr := strconv.Atoi(intAsString)
	// sanity.Ensure(parseErr != nil, "Invalid integer value [%s]", intAsString)
	return intValue
}

func (boolNode *Node) Bool() bool {
	sanity.Ensure(boolNode.Type == BOOL, "Invalid cast from type [%s] to [%s]", boolNode.Type, BOOL)
	boolAsString := boolNode.Values[0].(string)
	switch boolAsString {
	default:
		sanity.Fail("Invalid boolean value [%s]", boolAsString)
		return false
	case "true":
		return true
	case "false":
		return false
	}
}