package runtime

import (
	"math"
)

func doMath(lhs interface{}, rhs interface{}, operator string) interface{} {
	// TODO: Get to work for rationals and decimals too
	lhsValue := lhs.(int)
	rhsValue := rhs.(int)

	switch operator {
	default:
		return nil
	case "+":
		return lhsValue + rhsValue
	case "-":
		return lhsValue - rhsValue
	case "*":
		return lhsValue * rhsValue
	case "/":
		return lhsValue / rhsValue
	case "**":
		return int64(math.Pow(float64(lhsValue), float64(rhsValue)))
	case "%":
		return lhsValue % rhsValue
	case "<":
		return lhsValue < rhsValue
	case ">":
		return lhsValue > rhsValue
	case "<=":
		return lhsValue <= rhsValue
	case ">=":
		return lhsValue >= rhsValue
	case "==":
		return lhsValue == rhsValue
	}
}