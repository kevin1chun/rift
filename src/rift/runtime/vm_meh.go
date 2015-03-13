package runtime

import (
	"fmt"
	"strings"
	"rift/lang"
	"rift/support/collections"
	"rift/support/logging"
	"rift/support/sanity"
)

func mainRift(riftDefs []*lang.Node) *lang.Rift {
	for _, riftDef := range riftDefs {
		rift := riftDef.Rift()
		if rift.Name() == "main" {
			return rift
		}
	}

	return nil
}

func fillPredefs(env collections.PersistentMap) {
	env.Set("std:println", func(args []interface{}) interface{} {
		var stringedArgs []string
		for _, arg := range args {
			stringedArgs = append(stringedArgs, fmt.Sprintf("%v", arg))
		}
		fmt.Println(strings.Join(stringedArgs, ""))
		return nil
	})
}

func dereference(env collections.PersistentMap, ref *lang.Ref) interface{} {
	sanity.Ensure(env.Contains(ref.String()), "Undefined reference to [%s]", ref.String())
	return env.GetOrNil(ref.String())
}

func doAssignment(env collections.PersistentMap, assignment *lang.Assignment) interface{} {
	// TODO: Is lazy assignment okay here?
	env.Set(assignment.Ref().String(), assignment.Value())
	return nil
}

func doOperation(env collections.PersistentMap, op *lang.Operation) interface{} {
	// TODO: Right now must both be numeric
	lhsValue := evaluate(env, op.LHS()).(int)
	rhsValue := evaluate(env, op.RHS()).(int)
	switch op.Operator() {
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
		// TODO: Oops
		return lhsValue * rhsValue
	case "%":
		return lhsValue % rhsValue
	}
}

func evaluate(env collections.PersistentMap, a *lang.Node) interface{} {
	switch a.Type {
	default:
		return nil
	case lang.OP:
		return doOperation(env, a.Operation())
	case lang.ASSIGNMENT:
		return doAssignment(env, a.Assignment())
	case lang.FUNCAPPLY:
		return doFuncApply(env, a.FuncApply())
	case lang.REF:
		return evaluate(env, dereference(env, a.Ref()).(*lang.Node)) 
	case lang.STRING:
		return a.Str()
	case lang.NUM:
		return a.Int()
	case lang.BOOL:
		return a.Bool()
	}
}

func doFuncApply(env collections.PersistentMap, funcApply *lang.FuncApply) interface{} {
	f := dereference(env, funcApply.Ref()).(func([]interface{})interface{})
	args := funcApply.Args().Values()
	var argValues []interface{}
	for _, arg := range args {
		argValue := evaluate(env, arg.(*lang.Node))
		argValues = append(argValues, argValue)	
	}
	return f(argValues)
}

func Run(rifts []*lang.Node) {
	// TODO: Oops this only supports one rift :)
	if main := mainRift(rifts); main != nil {
		// ctx := collections.Stack{}
		env := collections.NewPersistentMap()
		fillPredefs(env)
		for _, line := range main.Lines() {
			evaluate(env, line)
		}
		// logging.Debug("Final environment: %+v", env.Freeze())
	} else {
		logging.Warn("No such rift [main]")
	}
}