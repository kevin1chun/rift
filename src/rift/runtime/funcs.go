package runtime

import (
	"rift/lang"
	"rift/support/collections"
	"rift/support/sanity"
)

func ensureArity(refStr string, expectedLength int, actualLength int) {
	sanity.Ensure(actualLength == expectedLength, "Function [%s] expects [%d] arguments, but got [%d]", refStr, expectedLength, actualLength)
}

func makeFunc(outerEnv collections.PersistentMap, f *lang.Func) func([]interface{}) interface{} {
	return func(args []interface{}) interface{} {
		// TODO: Assert arg list lengths match
		env := collections.ExtendPersistentMap(outerEnv)
		for i, argRef := range f.Args() {
			env.Set(argRef.String(), args[i])
		}
		
		var lastValue interface{}
		for _, line := range f.Lines() {
			lastValue = evaluate(env, line)
		}
		return lastValue
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
	returnValue := make(chan interface{}, 1)
	go func() {
		returnValue <- f(argValues)
	}()
	return <-returnValue
}
