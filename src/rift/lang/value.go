package lang

import (
	"bytes"
	"strconv"
	"rift/support/sanity"
)

func (stringNode *Node) Str() string {
	sanity.Ensure(stringNode.Type == STRING, "Invalid cast from type [%s] to [%s]", stringNode.Type, STRING)
	origStr := stringNode.Values[0].(string)

	var buffer bytes.Buffer
	var escaped bool
	for _, c := range origStr {
		if escaped {
			switch c {
			case '\\':
				buffer.WriteRune('\\')
			case '\'':
				buffer.WriteRune('\'')
			case '?':
				buffer.WriteRune('?')
			case '"':
				buffer.WriteRune('"')
			case 'a':
				buffer.WriteRune('\a')
			case 'b':
				buffer.WriteRune('\b')
			case 'f':
				buffer.WriteRune('\f')
			case 'n':
				buffer.WriteRune('\n')
			case 'r':
				buffer.WriteRune('\r')
			case 't':
				buffer.WriteRune('\t')
			case 'v':
				buffer.WriteRune('\v')
			}
			escaped = false
			continue
		}
		if c == '\\' {
			escaped = true
		} else {
			buffer.WriteRune(c)
		}
	}

	return buffer.String()
}

func (numericNode *Node) Int() int {
	sanity.Ensure(numericNode.Type == NUM, "Invalid cast from type [%s] to [%s]", numericNode.Type, NUM)
	intAsString := numericNode.Values[0].(string)
	intValue, parseErr := strconv.Atoi(intAsString)
	sanity.Ensure(parseErr == nil, "Invalid integer value [%s]", intAsString)
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
