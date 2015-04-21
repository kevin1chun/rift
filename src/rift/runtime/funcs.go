package runtime

import (
	"rift/lang"
	"rift/support/collections"
	"rift/support/sanity"
)

func ensureArity(expectedLength int, actualLength int) {
	sanity.Ensure(actualLength == expectedLength, "Function expects [%d] arguments, but got [%d]", expectedLength, actualLength)
}

func makeFunc(rift *lang.Rift, outerEnv collections.PersistentMap, f *lang.Func) func([]interface{}) interface{} {
	env := collections.ExtendPersistentMap(outerEnv)
	return func(args []interface{}) interface{} {
		ensureArity(len(f.Args()), len(args))
		for i, argRef := range f.Args() {
			env.Set(argRef.String(), args[i])
		}
		
		var lastValue interface{}
		for _, line := range f.Lines() {
			lastValue = evaluate(rift, env, line)
		}
		return lastValue
	}
}

func doFuncApply(rift *lang.Rift, env collections.PersistentMap, funcApply *lang.FuncApply) interface{} {
	f := dereference(rift, env, funcApply.Ref()).(func([]interface{})interface{})
	args := funcApply.Args().Values()
	var argValues []interface{}
	for _, arg := range args {
		argValue := evaluate(rift, env, arg.(*lang.Node))
		argValues = append(argValues, argValue)	
	}
	returnValue := make(chan interface{}, 1)
	go func() {
		returnValue <- f(argValues)
	}()
	return <-returnValue
}
