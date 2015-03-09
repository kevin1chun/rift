package lang

import (
	"rift/support/logging"
)

// const (
// 	RIFT  = "rift"
// 	FUNC = "function-definition"
// 	FUNCAPPLY = "function-apply"
// 	ARGS = "arguments"
// 	TUPLE = "tuple"
// 	LIST = "list"
// 	ASSIGNMENT = "assignment"
// 	IF = "if"
// 	STRING = "string"
// 	NUM = "numeric"
// 	BOOL = "boolean"
// 	REF = "reference"
// 	OP = "operation"
// 	BINOP = "binary-operator"
// )

const (
	
)

func mainRift(riftDefs []*Node) *Rift {
	for _, riftDef := range riftDefs {
		rift := riftDef.Rift()
		if rift.Name() == "main" {
			return rift
		}
	}

	return nil
}

func Compile(rifts []*Node) {
	if main := mainRift(rifts); main != nil {
		logging.Info("Found rift [main]: %s", main)
	} else {
		logging.Warn("No such rift [main]")
	}
}